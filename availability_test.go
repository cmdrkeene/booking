package booking

import (
	"database/sql"
	"os"
	"reflect"
	"testing"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

func TestAvailablityTable(t *testing.T) {
	path := "./availability_test.db"
	os.Remove(path)
	defer os.Remove(path)

	db, err := sql.Open("sqlite3", path)
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	table := newAvailabilityTable(db)

	feb1 := time.Date(2015, 1, 1, 0, 0, 0, 0, time.UTC)
	feb2 := feb1.Add(day)
	table.Add(feb1)
	table.Add(feb2)

	// list all
	list, err := table.List()
	if err != nil {
		t.Error(err)
	}

	// remove one
	want := []time.Time{feb1, feb2}
	if !reflect.DeepEqual(want, list) {
		t.Error("want", want)
		t.Error("got ", list)
	}

	// list to check removed
	table.Remove(feb2)
	want = []time.Time{feb1}
	list, err = table.List()
	if err != nil {
		t.Error(err)
	}
	if !reflect.DeepEqual(want, list) {
		t.Error("want", want)
		t.Error("got ", list)
	}
}
