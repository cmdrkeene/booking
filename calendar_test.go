package booking

import (
	"reflect"
	"testing"

	_ "github.com/mattn/go-sqlite3"
)

func TestCalendar(t *testing.T) {
	var cal Calendar
	cal.DB = testDB()
	defer cal.DB.Close()

	cal.Init()

	dates, err := cal.List()
	if err != nil {
		t.Error(err)
	}
	if l := len(dates); l != 0 {
		t.Error("want 0")
		t.Error("got ", l)
	}

	list := []date{
		newDate(2014, 1, 1),
		newDate(2014, 1, 2),
	}
	err = cal.Add(list...)
	if err != nil {
		t.Error(err)
	}

	dates, err = cal.List()
	if !reflect.DeepEqual(list, dates) {
		t.Error("want", list)
		t.Error("got ", dates)
	}

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
