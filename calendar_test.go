package booking

import (
	"reflect"
	"testing"

	"github.com/facebookgo/inject"
	_ "github.com/mattn/go-sqlite3"
)

func TestCalendar(t *testing.T) {
	db := testDB()
	defer db.Close()
	var cal Calendar
	err := inject.Populate(db, &cal)
	if err != nil {
		t.Error(err)
	}

	// List() -> []
	dates, err := cal.List()
	if err != nil {
		t.Error(err)
	}
	if l := len(dates); l != 0 {
		t.Error("want 0")
		t.Error("got ", l)
	}

	// Add() -> ok
	list := []date{
		newDate(2014, 1, 1),
		newDate(2014, 1, 2),
	}
	err = cal.Add(list...)
	if err != nil {
		t.Error(err)
	}

	// List() -> [2014-01-01 2014-01-02]
	dates, err = cal.List()
	if !reflect.DeepEqual(list, dates) {
		t.Error("want", list)
		t.Error("got ", dates)
	}

	// Available()
	var rangeTests = []struct {
		start date
		stop  date
		ok    bool
	}{
		{newDate(2014, 1, 1), newDate(2014, 1, 1), true},
		{newDate(2014, 1, 1), newDate(2014, 1, 2), true},
		{newDate(2014, 1, 1), newDate(2014, 1, 3), false},
		{newDate(2013, 12, 31), newDate(2014, 1, 2), false},
	}
	for _, tt := range rangeTests {
		if ok, _ := cal.Available(tt.start, tt.stop); ok != tt.ok {
			t.Log(tt.start, tt.stop)
			t.Error("want", tt.ok)
			t.Error("got ", ok)
		}
	}

	// Remove -> ok
	err = cal.Remove(list[0])
	if err != nil {
		t.Error(err)
	}

	dates, err = cal.List()
	list = []date{list[1]}
	if !reflect.DeepEqual(list, dates) {
		t.Error("want", list)
		t.Error("got ", dates)
	}
}
