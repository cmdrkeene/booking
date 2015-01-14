package booking

import (
	"database/sql"
	"time"
)

type availabilityStore interface {
	Add(time.Time) error
	Remove(time.Time) error
	List() ([]time.Time, error)
}

type availabilityTable struct {
	add    *sql.Stmt
	create *sql.Stmt
	remove *sql.Stmt
	list   *sql.Stmt
}

func newAvailabilityTable(db *sql.DB) availabilityTable {
	db.Exec(`
    CREATE TABLE Availability (
      Date DATETIME NOT NULL
    );
  `)

	t := availabilityTable{}
	var err error

	t.add, err = db.Prepare("INSERT INTO Availability (Date) VALUES ($1)")
	if err != nil {
		panic(err)
	}

	t.remove, err = db.Prepare("DELETE FROM Availability WHERE Date=$1")
	if err != nil {
		panic(err)
	}

	t.list, err = db.Prepare("SELECT Date FROM Availability ORDER BY Date ASC")
	if err != nil {
		panic(err)
	}

	return t
}

func (table availabilityTable) Add(t time.Time) error {
	_, err := table.add.Exec(t)
	return err
}

func (table availabilityTable) Remove(t time.Time) error {
	_, err := table.remove.Exec(t)
	return err
}

func (table availabilityTable) List() ([]time.Time, error) {
	rows, err := table.list.Query()
	if err != nil {
		return []time.Time{}, err
	}
	defer rows.Close()
	var list []time.Time
	for rows.Next() {
		var t time.Time
		err := rows.Scan(&t)
		if err != nil {
			return []time.Time{}, err
		}
		list = append(list, t)
	}
	err = rows.Err()
	if err != nil {
		return []time.Time{}, err
	}

	return list, nil
}
