package ta

import (
	"time"
)

type PeriodDefinition struct {
	size int8
	pad  int64
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

func (p periodDefinitions) calc(date time.Time) Period {
	var t Period
	for _, p := range p[:len(p)-1] {
		t += (p.cast(date) * Period(p.pad))
	}

	return t
}

/*

func (p periodDefinitions) extract(date time.Time) map[string]int64 {
    result := make(map[string]int64, 0)

    for _, p := range p[:len(p)-1] {
        result[p.name] =
    }

    fmt.Println("period", t)

    return t
}*/
