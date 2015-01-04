package booking

import (
	"testing"
	"time"
)

func TestDateRange(t *testing.T) {
	t1 := time.Date(2015, 2, 1, 0, 0, 0, 0, time.UTC)
	t2 := time.Date(2015, 2, 3, 0, 0, 0, 0, time.UTC)
	r := newDateRange(t1, t2)
	if r.Duration != 48*time.Hour {
		t.Error("want", 48*time.Hour)
		t.Error("got ", r.Duration)
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
