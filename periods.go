package aggregator

import "time"

type unit int

const (
	Year unit = 1 << iota
	Month
	Week
	Day
	YearDay
	Weekday
	Hour
	Minute
	Second
)

const binaryVersion uint64 = 1

type periodDefinition struct {
	size int16
	pad  uint64
	name string
	zero bool
	cast func(date time.Time) int
}

var units = []unit{Year, Month, Week, Day, Hour, Minute, Second, Weekday}
var defs = map[unit]periodDefinition{
	Year: {
		size: -1,
		pad:  1e4,
		name: "year",
		cast: func(d time.Time) int {
			return d.Year()
		},
	},
	Month: {
		size: 12,
		pad:  1e2,
		name: "month",
		zero: true,
		cast: func(d time.Time) int {
			return int(d.Month())
		},
	},
	Week: {
		size: 53,
		pad:  1e2,
		name: "week",
		zero: true,
		cast: func(d time.Time) int {
			w, _ := d.ISOWeek()
			return w
		},
	},
	Day: {
		size: 31,
		pad:  1e3,
		name: "day",
		cast: func(d time.Time) int {
			return d.Day()
		},
	},
	YearDay: {
		size: 367,
		pad:  1e3,
		name: "yearday",
		zero: true,
		cast: func(date time.Time) int {
			return date.YearDay()
		},
	},
	Weekday: {
		size: 7,
		pad:  1e1,
		name: "weekday",
		cast: func(date time.Time) int {
			return int(date.Weekday())
		},
	},
	Hour: {
		size: 24,
		pad:  1e2,
		name: "hour",
		cast: func(d time.Time) int {
			return d.Hour()
		},
	},
	Minute: {
		size: 60,
		pad:  1e2,
		name: "minute",
		cast: func(d time.Time) int {
			return d.Minute()
		},
	},
	Second: {
		size: 60,
		pad:  1e2,
		name: "second",
		cast: func(d time.Time) int {
			return d.Second()
		},
	},
}

type Period uint64

func newPeriod(flag unit, date time.Time) Period {
	us := getUnitsFromFlag(flag)

	t := binaryVersion
	for _, u := range us[:len(us)-1] {
		t *= defs[u].pad
		t += uint64(defs[u].cast(date))
	}

	t *= 1e3
	t += uint64(flag)

	return Period(t)
}

// ToMap returns a map representation of this period
func (p Period) ToMap() map[string]uint64 {
	t := uint64(p)

	us := p.Units()
	t = t / 1e3

	result := make(map[string]uint64, 0)
	for i := len(us) - 2; i >= 0; i-- {
		def := defs[us[i]]
		result[def.name] = t % def.pad
		t = t / def.pad
	}

	if t != binaryVersion {
		panic("Malformed period")
	}

	return result
}

func (p Period) flag() unit {
	return unit(uint64(p) % 1e3)
}

// Units returns a slice of the units of the period
func (p Period) Units() []unit {
	return getUnitsFromFlag(p.flag())
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
