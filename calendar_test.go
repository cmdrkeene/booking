package booking

import (
	"testing"
	"time"
)

func TestCalendar(t *testing.T) {
	var feb1 = time.Date(2015, 2, 1, 0, 0, 0, 0, time.UTC)
	var mar1 = time.Date(2015, 3, 1, 0, 0, 0, 0, time.UTC)

	var newCalendar = func() calendar {
		c := calendar{}
		c.SetAvailable(newDateRange(feb1, 2))
		c.SetAvailable(newDateRange(mar1, 2))
		c.SetBooked(newDateRange(mar1, 2))
		return c
	}

	var tests = []struct {
		dates   dateRange
		canBook bool
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
		if c.SetBooked(tt.dates) != tt.canBook {
			t.Error("requesting", tt.dates)
			t.Error("want", tt.canBook)
			t.Error("got ", !tt.canBook)
			t.Error(c)
		}
	}
}
