package booking

import (
	"database/sql"
	"os"
	"testing"

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

	list, err := table.List()
	if err != nil {
		t.Error(err)
	}

	if len(list) != 0 {
		t.Error("want", 0)
	}
}
