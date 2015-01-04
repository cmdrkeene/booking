package booking

import (
	"reflect"
	"testing"
	"time"
)

func TestDateRange(t *testing.T) {
	tests := []struct {
		t1 time.Time
		t2 time.Time
		r  dateRange
	}{
		// t1 before t2
		{
			time.Date(2015, 2, 1, 0, 0, 0, 0, time.UTC),
			time.Date(2015, 2, 3, 0, 0, 0, 0, time.UTC),
			dateRange{
				Start:    time.Date(2015, 2, 1, 0, 0, 0, 0, time.UTC),
				Duration: 48 * time.Hour,
			},
		},
		// t2 before t1
		{
			time.Date(2015, 2, 3, 0, 0, 0, 0, time.UTC),
			time.Date(2015, 2, 1, 0, 0, 0, 0, time.UTC),
			dateRange{
				Start:    time.Date(2015, 2, 1, 0, 0, 0, 0, time.UTC),
				Duration: 48 * time.Hour,
			},
		},
		// t1 equals t2
		{
			time.Date(2015, 2, 1, 0, 0, 0, 0, time.UTC),
			time.Date(2015, 2, 1, 0, 0, 0, 0, time.UTC),
			dateRange{
				Start:    time.Date(2015, 2, 1, 0, 0, 0, 0, time.UTC),
				Duration: 0,
			},
		},
	}

	for _, tt := range tests {
		r := newDateRange(tt.t1, tt.t2)
		if !reflect.DeepEqual(tt.r, r) {
			t.Error("want", tt.r)
			t.Error("got ", r)
		}
	}
}

func TestCalendar(t *testing.T) {
	// c := calendar{}
}

func TestBook(t *testing.T) {
	// success
	// fail unavailable
}

func TestGuestRegister(t *testing.T) {
	// success
	// fail existing user
	// fail missing name
	// fail missing email
	// fail short name
	// fail invalid email
}
