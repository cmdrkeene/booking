package booking

import (
	"reflect"
	"testing"
	"time"
)

func TestCalendar(t *testing.T) {
	var gId = guestId(999)
	var feb1 = time.Date(2015, 2, 1, 0, 0, 0, 0, time.UTC)
	var mar1 = time.Date(2015, 3, 1, 0, 0, 0, 0, time.UTC)

	var newTestCalendar = func() calendar {
		c := newCalendar()
		c.SetAvailable(newDateRange(feb1, 2))
		c.SetAvailable(newDateRange(mar1, 2))
		c.Reserve(newDateRange(mar1, 2), rateWithBunny, gId)
		return c
	}

	// Available()
	cal := newTestCalendar()
	available := []time.Time{
		feb1,
		feb1.Add(day),
	}
	if days := cal.Available(); !reflect.DeepEqual(available, days) {
		t.Error("want", available)
		t.Error("got ", days)
	}

	// Reserve()
	var tests = []struct {
		dates dateRange
		err   error
	}{
		{newDateRange(feb1, 1), nil},                   // available
		{newDateRange(feb1, 2), nil},                   // available
		{newDateRange(feb1, 3), unavailable},           // overran availability
		{newDateRange(feb1.Add(-day), 2), unavailable}, // underran availability
		{newDateRange(mar1, 1), unavailable},           // reserved
		{newDateRange(mar1, 2), unavailable},           // reserved
		{newDateRange(mar1.Add(day), 2), unavailable},  // reserved / overran availability
	}

	for _, tt := range tests {
		cal := newTestCalendar()
		if err := cal.Reserve(tt.dates, rateWithBunny, gId); err != tt.err {
			t.Error("requesting", tt.dates)
			t.Error("want", tt.err)
			t.Error("got ", err)
			t.Error(cal)
		}
	}
}
