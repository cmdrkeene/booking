package booking

import "testing"

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
