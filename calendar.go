package booking

import (
	"database/sql"

	"github.com/cmdrkeene/booking/pkg/date"
	"github.com/golang/glog"
)

// List of dates available for booking
// TODO turn Calendar into CalendarTx?
type Calendar struct {
	DB *sql.DB `inject:""`
}

type CalendarDB struct {
	*sql.DB
}

func (db *CalendarDB) Begin() (*CalendarTx, error) {
	tx, err := db.DB.Begin()
	if err != nil {
		return nil, err
	}
	return &CalendarTx{tx}, nil
}

type CalendarTx struct {
	*sql.Tx
}

func (tx *CalendarTx) Available(start, stop date.Date) (bool, error) {
	list, err := tx.List()
	if err != nil {
		return false, err
	}

	var include = func(l []date.Date, current date.Date) bool {
		for _, d := range l {
			if d == current {
				return true
			}
		}
		return false
	}

	daysApart := start.DaysApart(stop)
	for i := 0; i <= daysApart; i++ {
		current := start.Add(i)
		if !include(list, current) {
			return false, unavailable
		}
	}

	return true, nil
}

func (tx *CalendarTx) List() ([]date.Date, error) {
	stmt, err := tx.Prepare(`select Date from Calendar`)
	if err != nil {
		panic(err)
	}

	rows, err := stmt.Query()
	if err != nil {
		glog.Error(err)
		return []date.Date{}, err
	}
	defer rows.Close()

	var list []date.Date
	for rows.Next() {
		var d date.Date
		err := rows.Scan(&d)
		if err != nil {
			glog.Error(err)
			return []date.Date{}, err
		}
		list = append(list, d)
	}

	err = rows.Err()
	if err != nil {
		return []date.Date{}, err
	}

	return list, nil
}

const CalendarSchema = `
  CREATE TABLE Calendar (
    Date DATETIME UNIQUE NOT NULL
  )
`

func (c *Calendar) Add(dates ...date.Date) error {
	stmt, err := c.DB.Prepare(`insert into Calendar (Date) values ($1)`)
	if err != nil {
		panic(err)
	}

	for _, d := range dates {
		_, err := stmt.Exec(d)
		if err != nil {
			glog.Error(err)
			return err
		}
		glog.Infoln("added", d, "to availabilty calendar")
	}
	return nil
}

func (c *Calendar) Available(start, stop date.Date) (bool, error) {
	db := &CalendarDB{c.DB}
	tx, err := db.Begin()
	if err != nil {
		return false, err
	}

	ok, err := tx.Available(start, stop)
	if err != nil {
		tx.Rollback()
	}

	err = tx.Commit()
	if err != nil {
		glog.Error(err)
		return false, err
	}

	return ok, nil
}

func (c *Calendar) List() ([]date.Date, error) {
	db := &CalendarDB{c.DB}
	tx, err := db.Begin()
	if err != nil {
		return []date.Date{}, err
	}

	list, err := tx.List()
	if err != nil {
		tx.Rollback()
		glog.Error(err)
		return []date.Date{}, err
	}

	err = tx.Commit()
	if err != nil {
		glog.Error(err)
		return []date.Date{}, err
	}

	return list, nil
}

func (c *Calendar) Remove(dates ...date.Date) error {
	stmt, err := c.DB.Prepare(`delete from Calendar where Date=$1`)
	if err != nil {
		panic(err)
	}

	for _, d := range dates {
		_, err := stmt.Exec(d)
		if err != nil {
			glog.Error(err)
			return err
		}
		glog.Infoln("removed", d, "from availability calendar")
	}

	return nil
}
