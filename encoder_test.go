package aggregator

import (
	"encoding/hex"
	"encoding/json"

	"labix.org/v2/mgo/bson"

	. "gopkg.in/check.v1"
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
