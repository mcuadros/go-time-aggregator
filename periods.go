package ta

import (
	"time"
)

type PeriodDefinition struct {
	size int8
	pad  uint64
	name string
	cast func(date time.Time) Period
}

var Year = PeriodDefinition{
	size: -1,
	pad:  1e4,
	name: "year",
	cast: func(date time.Time) Period {
		return Period(date.Year())
	},
}

var Month = PeriodDefinition{
	size: 12,
	pad:  1e2,
	name: "month",
	cast: func(date time.Time) Period {
		return Period(date.Month())
	},
}

var Hour = PeriodDefinition{
	size: 24,
	pad:  1e2,
	name: "hour",
	cast: func(date time.Time) Period {
		return Period(date.Hour())
	},
}

var Weekday = PeriodDefinition{
	size: 7,
	pad:  1e1,
	name: "weekday",
	cast: func(date time.Time) Period {
		return Period(date.Weekday())
	},
}

type periodDefinitions []PeriodDefinition

func (p periodDefinitions) pack(date time.Time) Period {
	t := Period(1)
	for _, p := range p[:len(p)-1] {
		t *= Period(p.pad)
		t += p.cast(date)
	}

	return t
}

func (p periodDefinitions) unpack(total Period) map[string]uint64 {
	result := make(map[string]uint64, 0)
	t := uint64(total)
	for i := len(p) - 2; i >= 0; i-- {
		result[p[i].name] = t % p[i].pad
		t = t / p[i].pad
	}

	if t != 1 {
		panic("Malformed period")
	}

	return result
}
