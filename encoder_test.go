package aggregator

import (
	"bytes"
	"encoding/gob"
	"encoding/json"

	. "gopkg.in/check.v1"
	"gopkg.in/mgo.v2"
)

var example = "544140000000000000000100000000000000c94db70000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000a00000000000000"

func (s *UtilsSuite) Test_TimeAggregator_GetBSON(c *C) {
	a, _ := NewTimeAggregator(Year, Hour)
	a.Add(date2013December, 10)

	raw, err := a.GetBSON()
	c.Assert(raw.([]byte), HasLen, 218)
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

func (s *UtilsSuite) Test_TimeAggregator_MarshalJSON(c *C) {
	a, _ := NewTimeAggregator(Year, Hour)
	a.Add(date2015December, 10)
	a.Add(date2015November21h, 10)

	b, err := json.Marshal(a)
	c.Assert(err, IsNil)
	c.Assert(
		string(b),
		Equals,
		`[[{"year":2015},[0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,10,0,10]]]`,
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

func (s *UtilsSuite) TestIntegration(c *C) {
	session, err := mgo.Dial("localhost")
	db := session.DB("TimeAggregatorTest")
	defer db.DropDatabase()

	c.Assert(err, IsNil)
	collection := db.C("collection")

	a, _ := NewTimeAggregator(Year, Hour)
	a.Add(date2014November, 15)
	a.Add(date2015November, 10)
	a.Add(date2015December, 10)

	err = collection.Insert(struct{ A *TimeAggregator }{A: a})
	c.Assert(err, IsNil)

	r := struct{ A *TimeAggregator }{}
	err = collection.Find(nil).One(&r)
	c.Assert(err, IsNil)

	c.Assert(r.A.Values, HasLen, 2)
	c.Assert(r.A.Get(date2015January), Equals, int64(20))
	c.Assert(r.A.Get(date2014February), Equals, int64(15))
}
