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
	return c.withTx(
		func(tx *CalendarTx) error {
			return tx.Add(dates...)
		},
	)
}

func (c *Calendar) Available(start, stop date.Date) (bool, error) {
	var available bool
	err := c.withTx(
		func(tx *CalendarTx) error {
			var txErr error
			available, txErr = tx.Available(start, stop)
			return txErr
		},
	)
	return available, err
}

func (c *Calendar) List() ([]date.Date, error) {
	var dates []date.Date
	err := c.withTx(
		func(tx *CalendarTx) error {
			var txErr error
			dates, txErr = tx.List()
			return txErr
		},
	)
	return dates, err
}

func (c *Calendar) Remove(dates ...date.Date) error {
	return c.withTx(func(tx *CalendarTx) error {
		return tx.Remove(dates...)
	})
}

// Wraps fn to provide a CalendarTx that calls Rollback on err, Commit on ok
// Theoretically you could call a series of funcs for batch
func (c *Calendar) withTx(fn func(*CalendarTx) error) error {
	dbTx, err := c.DB.Begin()
	if err != nil {
		return err
	}
	tx := &CalendarTx{dbTx}
	err = fn(tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}

// CalendarTx wraps *sql.Tx to implement Calendar data operations
type CalendarTx struct {
	*sql.Tx
}

func (tx *CalendarTx) Add(dates ...date.Date) error {
	stmt, err := tx.Prepare(`insert into Calendar (Date) values ($1)`)
	if err != nil {
		return err
	}

	for _, d := range dates {
		_, err := stmt.Exec(d)
		if err != nil {
			return err
		}
		glog.Infoln("added", d, "to availabilty calendar")
	}
	return nil
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
		return []date.Date{}, err
	}

	rows, err := stmt.Query()
	if err != nil {
		return []date.Date{}, err
	}
	defer rows.Close()

	var list []date.Date
	for rows.Next() {
		var d date.Date
		err := rows.Scan(&d)
		if err != nil {

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

func (tx *CalendarTx) Remove(dates ...date.Date) error {
	stmt, err := tx.Prepare(`delete from Calendar where Date=$1`)
	if err != nil {
		return err
	}

	for _, d := range dates {
		_, err := stmt.Exec(d)
		if err != nil {
			return err
		}
		glog.Infoln("removed", d, "from availability calendar")
	}

	return nil
}
