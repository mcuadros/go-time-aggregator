package aggregator

import (
	"bytes"
	"encoding/binary"
	"errors"
	"io"
	"time"
)

var InvalidOrderError = errors.New("Invalid order of time units")

type TimeAggregator struct {
	Values map[Period]*aggregator
	kind   unit
	flags  unit
}

func NewTimeAggregator(units ...unit) (*TimeAggregator, error) {
	if !unitsAreSorted(units) {
		return nil, InvalidOrderError
	}

	return &TimeAggregator{
		Values: make(map[Period]*aggregator, 0),
		kind:   units[len(units)-1],
		flags:  unitsToFlag(units),
	}, nil
}

func (a *TimeAggregator) Add(date time.Time, value int64) {
	p := NewPeriod(a.flags, date)
	if _, ok := a.Values[p]; !ok {
		a.Values[p] = a.buildAggregator()
	}

	a.Values[p].Add(date, value)
}

func (a *TimeAggregator) Get(date time.Time) int64 {
	p := NewPeriod(a.flags, date)

	if _, ok := a.Values[p]; !ok {
		return -1
	}

	return a.Values[p].Get(date)
}

func (a *TimeAggregator) buildAggregator() *aggregator {
	return newAggregator(a.kind)
}

func (a *TimeAggregator) Marshal() []byte {
	buf := new(bytes.Buffer)

	for p, v := range a.Values {
		binary.Write(buf, binary.LittleEndian, p)
		buf.Write(v.Marshal())
	}

	return buf.Bytes()
}

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
			a.flags = p.Flag()
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

func newAggregator(u unit) *aggregator {
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

func unitsAreSorted(units []unit) bool {
	var p unit
	for _, u := range units {
		if p != 0 && u < p {
			return false
		}

		p = u
	}

	return true
}

func unitsToFlag(units []unit) unit {
	f := unit(0)
	for _, u := range units {
		f |= u
	}

	return f
}
