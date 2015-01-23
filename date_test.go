package booking

import (
	"reflect"
	"testing"
)

func TestNewDateFromString(t *testing.T) {
	var tests = []struct {
		input string
		err   error
		date  date
	}{
		{"", invalidDate, date{}},
		{"2015-01-02", nil, newDate(2015, 1, 2)},
	}

	for _, tt := range tests {
		date, err := newDateFromString(tt.input)
		if tt.err != err {
			t.Error("want", tt.err)
			t.Error("got ", err)
		}
		if !reflect.DeepEqual(tt.date, date) {
			t.Error("want", tt.date)
			t.Error("got ", date)
		}
	}
}

func TestDateDaysApart(t *testing.T) {
	var tests = []struct {
		d1 date
		d2 date
		n  int
	}{
		{newDate(2015, 1, 1), newDate(2015, 1, 1), 0},
		{newDate(2015, 1, 1), newDate(2015, 1, 2), 1},
		{newDate(2015, 1, 2), newDate(2015, 1, 1), 1},
	}
	for _, tt := range tests {
		if n := tt.d1.DaysApart(tt.d2); n != tt.n {
			t.Error("want", tt.n)
			t.Error("got ", n)
		}
	}
}
