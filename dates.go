package booking

import "time"

const day = 24 * time.Hour
const pretty = "January 2, 2006"
const iso8601 = "2006-01-02"

type dateRange struct {
	days  int // number of days from start, must be > 0
	start time.Time
}

func newDateRange(t time.Time, days int) dateRange {
	if days == 0 {
		panic("minimum days is 1")
	}
	return dateRange{start: t, days: days}
}

func (r dateRange) EachDay() []time.Time {
	var days []time.Time
	for i := 0; i < r.days; i++ {
		delta := time.Duration(i) * day
		days = append(days, r.start.Add(delta))
	}
	return days
}

func (r dateRange) String() string {
	t1 := r.start.Format(pretty)
	t2 := r.start.Add(time.Duration(r.days) * day).Format(pretty)
	return t1 + " to " + t2
}
