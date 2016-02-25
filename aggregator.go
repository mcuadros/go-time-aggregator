package aggregator

import (
	"bytes"
	"encoding/binary"
	"errors"
	"io"
	"time"
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
		flags:  unitsToFlag(units),
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
		return -1
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

func (a *TimeAggregator) buildAggregator() *aggregator {
	return newAggregator(a.kind)
}

// Marshal returns a binary representation of TimeAggregator
func (a *TimeAggregator) Marshal() []byte {
	buf := new(bytes.Buffer)

	for p, v := range a.Values {
		binary.Write(buf, binary.LittleEndian, p)
		buf.Write(v.Marshal())
	}

	return buf.Bytes()
}

// Unmarshal a binary string into a TimeAggregator, units and values are restored
func (a *TimeAggregator) Unmarshal(v []byte) error {
	a.Values = make(map[Period]*aggregator, 0)

	buf := bytes.NewBuffer(v)
	for {
		var p Period
		if err := binary.Read(buf, binary.LittleEndian, &p); err != nil {
			if err == io.EOF {
				break
			}

			return err
		}

		if a.flags == 0 {
			a.flags = p.flag()
			us := p.Units()

			a.kind = us[len(us)-1]
		}

		a.Values[p] = a.buildAggregator()
		if err := a.Values[p].Unmarshal(buf); err != nil {
			return err
		}
	}

	return nil
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
