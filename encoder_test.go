package aggregator

import (
	"bytes"
	"encoding/gob"
	"encoding/hex"
	"encoding/json"

	. "gopkg.in/check.v1"
	"gopkg.in/mgo.v2/bson"
)

var example = "c800000000094eb70000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000001400000000000000"

func (s *UtilsSuite) Test_TimeAggregator_GetBSON(c *C) {
	a, _ := NewTimeAggregator(Year, Hour)
	a.Add(date2013December, 10)

	raw, err := a.GetBSON()

	c.Assert(raw.([]byte), HasLen, 200)
	c.Assert(err, IsNil)
}

func (s *UtilsSuite) Test_TimeAggregator_GetBSON_nil(c *C) {
	var a *TimeAggregator
	raw, err := a.GetBSON()

	c.Assert(raw, IsNil)
	c.Assert(err, IsNil)
}

func (s *UtilsSuite) Test_TimeAggregator_GetBSON_empty(c *C) {
	a, _ := NewTimeAggregator(Month)
	raw, err := a.GetBSON()

	c.Assert(raw, IsNil)
	c.Assert(err, IsNil)
}

func (s *UtilsSuite) Test_TimeAggregator_SetBSON(c *C) {
	a, _ := NewTimeAggregator(Year, Hour)
	a.Add(date2013December, 10)

	raw := bson.Raw{}
	raw.Kind = 5
	raw.Data, _ = hex.DecodeString(example)

	b := &TimeAggregator{}
	err := b.SetBSON(raw)
	c.Assert(err, IsNil)
	c.Assert(b.Values, HasLen, 1)
}

func (s *UtilsSuite) Test_TimeAggregator_MarshalJSON(c *C) {
	a, _ := NewTimeAggregator(Year, Hour)
	a.Add(date2013December, 10)

	b, err := json.Marshal(a)
	c.Assert(err, IsNil)
	c.Assert(
		string(b),
		Equals,
		`[[{"year":2013},[0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,10]]]`,
	)
}

func (s *UtilsSuite) Test_TimeAggregator_Gob(c *C) {
	ta1, _ := NewTimeAggregator(Year, Hour)
	ta1.Add(date2013December, 10)

	var (
		buf bytes.Buffer
		err error
	)

	enc := gob.NewEncoder(&buf)
	err = enc.Encode(ta1)
	c.Assert(err, IsNil)

	var ta2 *TimeAggregator
	dec := gob.NewDecoder(&buf)
	err = dec.Decode(&ta2)
	c.Assert(err, IsNil)

	// Can't use `c.Assert(ta1, DeepEquals, ta2)` because it errors out.
	// Comparing closures is not straightforward.
	for p, a := range ta2.Values {
		c.Assert(a.values, DeepEquals, ta1.Values[p].values)
		c.Assert(a.p.size, DeepEquals, ta1.Values[p].p.size)
		c.Assert(a.p.pad, DeepEquals, ta1.Values[p].p.pad)
		c.Assert(a.p.name, DeepEquals, ta1.Values[p].p.name)
		c.Assert(a.p.zero, DeepEquals, ta1.Values[p].p.zero)
	}
}
