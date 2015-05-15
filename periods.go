package ta

import (
	"time"
)

type Period uint64

type unit int

const (
	Year unit = 1 << iota
	Month
	Week
	Day
	Hour
	Minute
	Second
	Weekday
)

type PeriodDefinition struct {
	size int8
	pad  uint64
	name string
	cast func(date time.Time) Period
}

var units = []unit{Year, Month, Week, Day, Hour, Minute, Second, Weekday}
var defs = map[unit]PeriodDefinition{
	Year: PeriodDefinition{
		size: -1,
		pad:  1e4,
		name: "year",
		cast: func(date time.Time) Period {
			return Period(date.Year())
		},
	},
	Month: PeriodDefinition{
		size: 12,
		pad:  1e2,
		name: "month",
		cast: func(date time.Time) Period {
			return Period(date.Month())
		},
	},
	Hour: PeriodDefinition{
		size: 24,
		pad:  1e2,
		name: "hour",
		cast: func(date time.Time) Period {
			return Period(date.Hour())
		},
	},
	Weekday: PeriodDefinition{
		size: 7,
		pad:  1e1,
		name: "weekday",
		cast: func(date time.Time) Period {
			return Period(date.Weekday())
		},
	},
}

func pack(flag unit, date time.Time) Period {
	us := getUnitsFromFlag(flag)

	t := Period(1)
	for _, u := range us[:len(us)-1] {
		t *= Period(defs[u].pad)
		t += defs[u].cast(date)
	}

	t *= 1e3
	t += Period(flag)

	return t
}

func unpack(total Period) map[string]uint64 {
	t := uint64(total)

	us := getUnitsFromFlag(unit(t % 1e3))
	t = t / 1e3

	result := make(map[string]uint64, 0)
	for i := len(us) - 2; i >= 0; i-- {
		def := defs[us[i]]
		result[def.name] = t % def.pad
		t = t / def.pad
	}

	if t != 1 {
		panic("Malformed period")
	}

	return result
}

func getUnitsFromFlag(flag unit) []unit {
	us := make([]unit, 0)
	for _, u := range units {
		if flag&u != 0 {
			us = append(us, u)
		}
	}

	return us
}
