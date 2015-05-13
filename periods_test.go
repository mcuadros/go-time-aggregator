package ta

import (
	. "gopkg.in/check.v1"
)

func (s *UtilsSuite) TestperiodDefinitions_pack(c *C) {
	ps := periodDefinitions{Year, Month, Weekday}
	c.Assert(ps.pack(date2015November), Equals, Period(1201511))
}

func (s *UtilsSuite) TestperiodDefinitions_unpack(c *C) {
	ps := periodDefinitions{Year, Month, Weekday}
	p := ps.unpack(Period(1201511))
	c.Assert(p, DeepEquals, map[string]uint64{"year": 2015, "month": 11})
}
