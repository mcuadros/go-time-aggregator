# go-time-aggregator [![GoDoc](https://godoc.org/github.com/mcuadros/go-time-aggregator?status.svg)](https://godoc.org/github.com/mcuadros/go-time-aggregator) [![Build Status](https://travis-ci.org/mcuadros/go-time-aggregator.svg)](https://travis-ci.org/mcuadros/go-time-aggregator)

Examples
--------

Adding to an aggrergator the number of hours that an employee works in 2016~2017 working days but August 
```
yearLog, _ := aggregator.NewTimeAggregator(aggregator.Year, aggregator.YearDay)
startingDay := time.Date(2016, 1, 1, 0, 0, 0, 1, time.UTC)
for i := 0; i < 365*2; i++ {
    day := startingDay.Add(time.Duration(i*24) * time.Hour)
    if day.Weekday() > 0 && day.Weekday() < 6 && day.Month() != 8 {
        yearLog.Add(day, int64(8))
    }
}
```

Printing a representation of the days worked
```
fmt.Println(yearLog))
```

Output:
```
Year: 2016	▁█▁▁█████▁▁█████▁▁█████▁▁█████▁▁█████▁▁█████▁▁█████▁▁█████▁▁█████▁▁█████▁▁█████▁▁█████▁▁█████▁▁█████▁▁█████▁▁█████▁▁█████▁▁█████▁▁█████▁▁█████▁▁█████▁▁█████▁▁█████▁▁█████▁▁█████▁▁█████▁▁█████▁▁█████▁▁█████▁▁█████▁▁▁▁▁▁▁▁▁▁▁▁▁▁▁▁▁▁▁▁▁▁▁▁▁▁▁▁▁▁▁▁▁██▁▁█████▁▁█████▁▁█████▁▁█████▁▁█████▁▁█████▁▁█████▁▁█████▁▁█████▁▁█████▁▁█████▁▁█████▁▁█████▁▁█████▁▁█████▁▁█████▁▁█████▁  yeardays
Year: 2017	▁▁█████▁▁█████▁▁█████▁▁█████▁▁█████▁▁█████▁▁█████▁▁█████▁▁█████▁▁█████▁▁█████▁▁█████▁▁█████▁▁█████▁▁█████▁▁█████▁▁█████▁▁█████▁▁█████▁▁█████▁▁█████▁▁█████▁▁█████▁▁█████▁▁█████▁▁█████▁▁█████▁▁█████▁▁█████▁▁█████▁▁█▁▁▁▁▁▁▁▁▁▁▁▁▁▁▁▁▁▁▁▁▁▁▁▁▁▁▁▁▁▁▁█▁▁█████▁▁█████▁▁█████▁▁█████▁▁█████▁▁█████▁▁█████▁▁█████▁▁█████▁▁█████▁▁█████▁▁█████▁▁█████▁▁█████▁▁█████▁▁█████▁▁█████▁▁▁  yeardays
```

License
-------

MIT, see [LICENSE](LICENSE)