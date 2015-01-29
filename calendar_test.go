package booking

import (
	"reflect"
	"testing"

	"github.com/cmdrkeene/booking/pkg/date"
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
	list := []date.Date{
		date.New(2014, 1, 1),
		date.New(2014, 1, 2),
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
		start date.Date
		stop  date.Date
		ok    bool
	}{
		{date.New(2014, 1, 1), date.New(2014, 1, 1), true},
		{date.New(2014, 1, 1), date.New(2014, 1, 2), true},
		{date.New(2014, 1, 1), date.New(2014, 1, 3), false},
		{date.New(2013, 12, 31), date.New(2014, 1, 2), false},
	}
	tx, err := db.Begin()
	if err != nil {
		t.Fatal(err)
	}
	for _, tt := range rangeTests {
		if ok, _ := cal.Available(tx, tt.start, tt.stop); ok != tt.ok {
			t.Log(tt.start, tt.stop)
			t.Error("want", tt.ok)
			t.Error("got ", ok)
		}
	}
	tx.Commit()

	// Remove -> ok
	err = cal.Remove(list[0])
	if err != nil {
		t.Error(err)
	}

	dates, err = cal.List()
	list = []date.Date{list[1]}
	if !reflect.DeepEqual(list, dates) {
		t.Error("want", list)
		t.Error("got ", dates)
	}
}
