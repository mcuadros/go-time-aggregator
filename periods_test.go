package aggregator

import . "gopkg.in/check.v1"

func (s *TimeAggregatorSuite) TestNewPeriod(c *C) {
	p := newPeriod(Year|Month, date2015November)
	c.Assert(p, Equals, Period(1201511003))

	p = newPeriod(Year|YearDay, date2015November)
	c.Assert(p, Equals, Period(12015316017))
}

func (s *TimeAggregatorSuite) TestPeriodToMap(c *C) {
	p := Period(1201511003)
	m, err := p.Map()
	c.Assert(err, IsNil)

	c.Assert(m, DeepEquals, map[string]uint64{"year": 2015, "month": 11})

	p = Period(12015316017)
	m, err = p.Map()
	c.Assert(err, IsNil)
	c.Assert(m, DeepEquals, map[string]uint64{"year": 2015, "yearday": 316})
}

func (s *TimeAggregatorSuite) TestPeriodString(c *C) {
	p := Period(1201511003)
	c.Assert(p.String(), Equals, "Year: 2015 / Month: 11")
}
