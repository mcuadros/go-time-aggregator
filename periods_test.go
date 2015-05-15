package ta

import (
	. "gopkg.in/check.v1"
)

func (s *UtilsSuite) Test_pack(c *C) {
	flags := Year | Month | Hour
	c.Assert(pack(flags, date2015November), Equals, Period(1201511019))
}

func (s *UtilsSuite) Test_unpack(c *C) {
	p := unpack(Period(1201511019))
	c.Assert(p, DeepEquals, map[string]uint64{"year": 2015, "month": 11})
}
