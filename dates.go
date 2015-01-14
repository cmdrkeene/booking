package booking

import "time"

const day = 24 * time.Hour
const pretty = "January 2, 2006"
const iso8601 = "2006-01-02"

type dateRange struct {
	list []time.Time
}

func newDateRange(start time.Time, numDays int) dateRange {
	if numDays == 0 {
		panic("numDays must be > 0")
	}

	dates := dateRange{}
	for i := 0; i < numDays; i++ {
		delta := time.Duration(i) * day
		dates.list = append(dates.list, start.Add(delta))
	}

	return dates
}

func newDateRangeBetween(start, end time.Time) dateRange {
	numDays := int(end.Sub(start).Hours() / 24)
	return newDateRange(start, numDays)
}

func (r dateRange) Start() time.Time {
	return r.list[0]
}

func (r dateRange) End() time.Time {
	return r.list[len(r.list)-1]
}

// return true if all times in range present in list
func (r dateRange) Coincident(list []time.Time) bool {
	set := make(map[time.Time]interface{})
	for _, t := range list {
		set[t] = struct{}{}
	}
	for _, t := range r.list {
		if _, ok := set[t]; !ok {
			return false
		}
	}
	return true
}

func (r dateRange) EachDay() []time.Time {
	return r.list
}

func (r dateRange) String() string {
	return r.Start().Format(pretty) + " to " + r.End().Format(pretty)
}
