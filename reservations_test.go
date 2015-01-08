package booking

import (
	"testing"
	"time"
)

func TestCalendar(t *testing.T) {
	var resId = reservationId("999")
	var feb1 = time.Date(2015, 2, 1, 0, 0, 0, 0, time.UTC)
	var mar1 = time.Date(2015, 3, 1, 0, 0, 0, 0, time.UTC)

	var newCalendar = func() calendar {
		c := calendar{}
		c.SetAvailable(newDateRange(feb1, 2))
		c.SetAvailable(newDateRange(mar1, 2))
		c.Reserve(newDateRange(mar1, 2), resId)
		return c
	}

	var tests = []struct {
		dates      dateRange
		canReserve bool
	}{
		{newDateRange(feb1, 1), true},            // available
		{newDateRange(feb1, 2), true},            // available
		{newDateRange(feb1, 3), false},           // overran availability
		{newDateRange(feb1.Add(-day), 2), false}, // underran availability
		{newDateRange(mar1, 1), false},           // booked
		{newDateRange(mar1, 2), false},           // booked
		{newDateRange(mar1.Add(day), 2), false},  // booked / overran availability
	}

	for _, tt := range tests {
		c := newCalendar()
		if c.Reserve(tt.dates, resId) != tt.canReserve {
			t.Error("requesting", tt.dates)
			t.Error("want", tt.canReserve)
			t.Error("got ", !tt.canReserve)
			t.Error(c)
		}
	}
}
