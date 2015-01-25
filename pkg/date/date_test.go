package date

import (
	"reflect"
	"testing"
	"time"
)

func TestFormat(t *testing.T) {
	var tests = []struct {
		date   Date
		layout string
		output string
	}{
		{New(2015, 2, 14), ISO8601, "2015-02-14"},
		{New(2015, 2, 14), Pretty, "February 14, 2015"},
	}

	for _, tt := range tests {
		if s := tt.date.Format(tt.layout); s != tt.output {
			t.Error("want", tt.output)
			t.Error("got ", s)
		}
	}
}

func TestParse(t *testing.T) {
	var tests = []struct {
		src  interface{}
		date Date
		err  error
	}{
		{time.Date(2015, 2, 3, 4, 5, 6, 7, time.UTC), New(2015, 2, 3), nil},
		{"", Date{}, ParseError},
		{"2015-01-02", New(2015, 1, 2), nil},
	}

	for _, tt := range tests {
		date, err := Parse(tt.src)
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

func TestDaysApart(t *testing.T) {
	var tests = []struct {
		d1 Date
		d2 Date
		n  int
	}{
		{New(2015, 1, 1), New(2015, 1, 1), 0},
		{New(2015, 1, 1), New(2015, 1, 2), 1},
		{New(2015, 1, 2), New(2015, 1, 1), 1},
		{New(2015, 1, 1), New(2015, 1, 5), 4},
	}
	for _, tt := range tests {
		if n := tt.d1.DaysApart(tt.d2); n != tt.n {
			t.Error("want", tt.n)
			t.Error("got ", n)
		}
	}
}
