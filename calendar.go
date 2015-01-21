package booking

import (
	"database/sql"

	"github.com/golang/glog"
)

// List of dates available for booking
type Calendar struct {
	DB     *sql.DB `inject:""`
	add    *sql.Stmt
	list   *sql.Stmt
	remove *sql.Stmt
}

// Init creates table if needed
func (c *Calendar) Init() {
	// create table
	_, err := c.DB.Exec(`
    create table calendar (
      date datetime
    )
  `)
	if err == nil {
		glog.Info("calendar table created")
	} else {
		glog.Warning(err)
	}

	// setup prepated statements
	c.add, err = c.DB.Prepare(`insert into calendar (date) values ($1)`)
	if err != nil {
		panic(err)
	}

	c.list, err = c.DB.Prepare(`select date from calendar`)
	if err != nil {
		panic(err)
	}

	c.remove, err = c.DB.Prepare(`delete from calendar where date=$1`)
	if err != nil {
		panic(err)
	}
}

func (c *Calendar) Add(dates ...date) error {
	for _, d := range dates {
		_, err := c.add.Exec(d)
		if err != nil {
			return err
		}
	}
	return nil
}

func (c *Calendar) Remove(dates ...date) error {
	for _, d := range dates {
		_, err := c.remove.Exec(d)
		if err != nil {
			return err
		}
	}
	return nil
}

func (c *Calendar) List() ([]date, error) {
	rows, err := c.list.Query()
	if err != nil {
		return []date{}, err
	}
	defer rows.Close()

	var list []date
	for rows.Next() {
		var d date
		err := rows.Scan(&d)
		if err != nil {
			return []date{}, err
		}
		list = append(list, d)
	}

	return list, nil
}
