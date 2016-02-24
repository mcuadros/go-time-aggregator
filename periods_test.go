package aggregator

import . "gopkg.in/check.v1"

func (s *UtilsSuite) Test_NewPeriod(c *C) {
	p := newPeriod(Year|Month, date2015November)
	c.Assert(p, Equals, Period(1201511003))

	p = newPeriod(Year|YearDay, date2015November)
	c.Assert(p, Equals, Period(12015316017))
}

func (s *UtilsSuite) Test_unpack(c *C) {
	p := Period(1201511003)
	c.Assert(p.ToMap(), DeepEquals, map[string]uint64{"year": 2015, "month": 11})

	p = Period(12015316017)
	c.Assert(p.ToMap(), DeepEquals, map[string]uint64{"year": 2015, "yearday": 316})
}
