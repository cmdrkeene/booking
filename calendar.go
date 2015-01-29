package booking

import (
	"database/sql"

	"github.com/cmdrkeene/booking/pkg/date"
	"github.com/golang/glog"
)

// List of dates available for booking
type Calendar struct {
	DB *sql.DB `inject:""`
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

func (c *Calendar) Available(tx *sql.Tx, start, stop date.Date) (bool, error) {
	list, err := c.ListTx(tx)
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

func (c *Calendar) List() ([]date.Date, error) {
	tx, err := c.DB.Begin()
	if err != nil {
		return []date.Date{}, err
	}

	list, err := c.ListTx(tx)
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

func (c *Calendar) ListTx(tx *sql.Tx) ([]date.Date, error) {
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
