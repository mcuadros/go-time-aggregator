package aggregator

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"sort"
	"time"

	"github.com/joliv/spark"
)

var (
	InvalidOrderError = errors.New("Invalid order of time units")
	UnitsMismatch     = errors.New("Units mismatch")
)

type TimeAggregator struct {
	Values map[Period]*aggregator
	kind   Unit
	flags  Unit
}

// NewTimeAggregator returns a new TimeAggregator configured with the given
// units, InvalidOrderError is return if the units are not sorted
func NewTimeAggregator(units ...Unit) (*TimeAggregator, error) {
	if !unitsAreSorted(units) {
		return nil, InvalidOrderError
	}

	return &TimeAggregator{
		Values: make(map[Period]*aggregator, 0),
		kind:   units[len(units)-1],
		flags:  unitsToFlag(units[:len(units)-1]),
	}, nil
}

// Add sum a value to the correct count using date
func (a *TimeAggregator) Add(date time.Time, value int64) {
	p := newPeriod(a.flags, date)
	if _, ok := a.Values[p]; !ok {
		a.Values[p] = a.buildAggregator()
	}

	a.Values[p].Add(date, value)
}

// Get returns the count for the given date
func (a *TimeAggregator) Get(date time.Time) int64 {
	p := newPeriod(a.flags, date)

	if _, ok := a.Values[p]; !ok {
		return 0
	}

	return a.Values[p].Get(date)
}

// Sum sum a TimeAggregator to other, if the units are diferent UnitsMismatch
// error is returned
func (a *TimeAggregator) Sum(b *TimeAggregator) error {
	if a.flags != b.flags {
		return UnitsMismatch
	}

	for p, v := range b.Values {
		if _, ok := a.Values[p]; !ok {
			a.Values[p] = v
		} else {
			a.Values[p].Sum(v)
		}
	}

	return nil
}

type Entry struct {
	Period map[string]uint64
	Value  int64
}

// Entries return a list of Entry structs with non-zero values
func (a *TimeAggregator) Entries() []Entry {
	var o []Entry

	var periods Periods
	for i := range a.Values {
		periods = append(periods, i)
	}

	sort.Sort(periods)
	for _, p := range periods {
		o = append(o, a.Values[p].entries(p)...)
	}

	return o
}

func (a *TimeAggregator) buildAggregator() *aggregator {
	return newAggregator(a.kind)
}

var signature = []byte{'T', 'A'}

// Marshal returns a binary representation of TimeAggregator
func (a *TimeAggregator) Marshal() []byte {
	buf := new(bytes.Buffer)
	buf.Write(signature)

	binary.Write(buf, binary.LittleEndian, int64(a.kind))
	binary.Write(buf, binary.LittleEndian, int64(a.flags))

	for p, v := range a.Values {
		binary.Write(buf, binary.LittleEndian, p)
		buf.Write(v.Marshal())
	}

	return buf.Bytes()
}

// Unmarshal a binary string into a TimeAggregator, units and values are restored
func (a *TimeAggregator) Unmarshal(v []byte) error {
	buf := bytes.NewBuffer(v)

	s := make([]byte, len(signature))
	if _, err := buf.Read(s); err != nil {
		return err
	}

	if string(s) != string(signature) {
		return fmt.Errorf("Signature missmatch found %q", s)
	}

	var kind Unit
	if err := binary.Read(buf, binary.LittleEndian, &kind); err != nil {
		return err
	}

	var flags Unit
	if err := binary.Read(buf, binary.LittleEndian, &flags); err != nil {
		return err
	}

	a.Values = make(map[Period]*aggregator, 0)
	a.kind = kind
	a.flags = flags

	for {
		var p Period
		if err := binary.Read(buf, binary.LittleEndian, &p); err != nil {
			if err == io.EOF {
				break
			}

			return err
		}

		a.Values[p] = a.buildAggregator()
		if err := a.Values[p].Unmarshal(buf); err != nil {
			return err
		}
	}

	return nil
}

func (a *TimeAggregator) String() string {
	var periods Periods
	for i := range a.Values {
		periods = append(periods, i)
	}

	sort.Sort(periods)

	var o string
	for _, p := range periods {
		o += fmt.Sprintf("%s\t%s  %ss\n", p, a.Values[p], a.Values[p].p.name)
	}

	return o

}

type aggregator struct {
	values []int64
	p      periodDefinition
}

func newAggregator(u Unit) *aggregator {
	return &aggregator{
		values: make([]int64, defs[u].size),
		p:      defs[u],
	}
}

func (a *aggregator) Add(date time.Time, value int64) {
	key := a.p.cast(date)
	if a.p.zero {
		key--
	}

	a.values[key] += value
}

func (a *aggregator) Get(date time.Time) int64 {
	key := a.p.cast(date)
	if a.p.zero {
		key--
	}

	return a.values[key]
}

func (a *aggregator) Sum(b *aggregator) {
	for i, value := range b.values {
		a.values[i] += value
	}
}

func (a *aggregator) entries(p Period) []Entry {
	var o []Entry
	m := p.Map()

	for i, v := range a.values {
		if v == 0 {
			continue
		}

		e := Entry{Period: m, Value: v}
		e.Period[a.p.name] = uint64(i)

		o = append(o, e)
	}

	return o
}

func (a *aggregator) Marshal() []byte {
	buf := new(bytes.Buffer)

	for _, v := range a.values {
		binary.Write(buf, binary.LittleEndian, v)
	}

	return buf.Bytes()
}

func (a *aggregator) Unmarshal(r io.Reader) error {
	for i, _ := range a.values {
		err := binary.Read(r, binary.LittleEndian, &a.values[i])
		if err != nil {
			return err
		}
	}

	return nil
}

func (a *aggregator) String() string {
	values := make([]float64, len(a.values))
	for i, v := range a.values {
		values[i] = float64(v)
	}

	return spark.Line(values)
}

func unitsAreSorted(units []Unit) bool {
	var p Unit
	for _, u := range units {
		if p != 0 && u < p {
			return false
		}

		p = u
	}

	return true
}

func unitsToFlag(units []Unit) Unit {
	f := Unit(0)
	for _, u := range units {
		f |= u
	}

	return f
}
