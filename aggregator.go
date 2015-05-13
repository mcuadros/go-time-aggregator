package ta

import (
	"bytes"
	"encoding/binary"
	"io"
	"time"
)

type Period int64

type TimeAggregator struct {
	Values  map[Period]*Aggregator
	periods periodDefinitions
}

func NewTimeAggregator(p ...PeriodDefinition) *TimeAggregator {
	return &TimeAggregator{
		Values:  make(map[Period]*Aggregator, 0),
		periods: periodDefinitions(p),
	}
}

func (a *TimeAggregator) Add(date time.Time, value int64) {
	p := a.periods.calc(date)

	if _, ok := a.Values[p]; !ok {
		a.Values[p] = a.buildAggregator()
	}

	a.Values[p].Add(date, value)
}

func (a *TimeAggregator) Get(date time.Time) int64 {
	p := a.periods.calc(date)

	if _, ok := a.Values[p]; !ok {
		return -1
	}

	return a.Values[p].Get(date)
}

func (a *TimeAggregator) buildAggregator() *Aggregator {
	return NewAggregator(a.periods[len(a.periods)-1])
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
	buf := bytes.NewBuffer(v)

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

type Aggregator struct {
	values []int64
	p      PeriodDefinition
}

func NewAggregator(p PeriodDefinition) *Aggregator {
	return &Aggregator{
		values: make([]int64, p.size),
		p:      p,
	}
}

func (a *Aggregator) Add(date time.Time, value int64) {
	a.values[a.p.cast(date)] += value
}

func (a *Aggregator) Get(date time.Time) int64 {
	return a.values[a.p.cast(date)]
}

func (a *Aggregator) Marshal() []byte {
	buf := new(bytes.Buffer)

	for _, v := range a.values {
		binary.Write(buf, binary.LittleEndian, v)
	}

	return buf.Bytes()
}

func (a *Aggregator) Unmarshal(r io.Reader) error {
	for i, _ := range a.values {
		err := binary.Read(r, binary.LittleEndian, &a.values[i])
		if err != nil {
			return err
		}
	}

	return nil
}
