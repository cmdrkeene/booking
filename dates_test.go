package booking

import (
	"reflect"
	"testing"
	"time"
)

func TestDateRange(t *testing.T) {
	// newDateRange
	feb1 := time.Date(2015, 2, 1, 0, 0, 0, 0, time.UTC)
	est, err := time.LoadLocation("America/New_York")
	if err != nil {
		t.Error(nil)
	}
	feb1Local := time.Date(2015, 1, 31, 19, 0, 0, 0, est)
	var tests = []struct {
		input   time.Time
		numDays int
		start   time.Time
		end     time.Time
		list    []time.Time
	}{
		{feb1, 1, feb1, feb1, []time.Time{feb1}},
		{feb1Local, 1, feb1, feb1, []time.Time{feb1}},
		{feb1, 2, feb1, feb1.Add(day), []time.Time{feb1, feb1.Add(day)}},
		{feb1, 3, feb1, feb1.Add(2 * day), []time.Time{feb1, feb1.Add(day), feb1.Add(2 * day)}},
	}
	for _, tt := range tests {
		dates := newDateRange(tt.input, tt.numDays)
		if dates.Start() != tt.start {
			t.Error("want", tt.start)
			t.Error("got ", dates.Start())
		}
		if dates.End() != tt.end {
			t.Error("want", tt.end)
			t.Error("got ", dates.End())
		}
		if !reflect.DeepEqual(dates.EachDay(), tt.list) {
			t.Error("want", tt.list)
			t.Error("got ", dates.EachDay())
		}
	}

	// newDateRangeBetween
	var betweenTests = []struct {
		t1    time.Time
		t2    time.Time
		start time.Time
		end   time.Time
		list  []time.Time
	}{
		{feb1, feb1, feb1, feb1, []time.Time{feb1}},
		{feb1Local, feb1Local, feb1, feb1, []time.Time{feb1}},
		{feb1, feb1.Add(day), feb1, feb1.Add(day), []time.Time{feb1, feb1.Add(day)}},
	}
	for _, tt := range betweenTests {
		dates := newDateRangeBetween(tt.t1, tt.t2)
		if dates.Start() != tt.start {
			t.Error("want", tt.start)
			t.Error("got ", dates.Start())
		}
		if dates.End() != tt.end {
			t.Error("want", tt.end)
			t.Error("got ", dates.End())
		}
		if !reflect.DeepEqual(dates.EachDay(), tt.list) {
			t.Error("want", tt.list)
			t.Error("got ", dates.EachDay())
		}
	}
}
