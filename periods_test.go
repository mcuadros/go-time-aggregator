package aggregator

import (
	. "gopkg.in/check.v1"
)

func (s *UtilsSuite) Test_NewPeriod(c *C) {
	p := newPeriod(Year|Month|Hour, date2015November)
	c.Assert(p, Equals, Period(1201511067))
}

func (s *UtilsSuite) Test_unpack(c *C) {
	p := Period(1201511067)
	c.Assert(p.ToMap(), DeepEquals, map[string]uint64{"year": 2015, "month": 11})
}
