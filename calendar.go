package booking

import (
	"database/sql"

	"github.com/golang/glog"
)

// List of dates available for booking
type Calendar struct {
	DB           *sql.DB `inject:""`
	add          *sql.Stmt
	list         *sql.Stmt
	remove       *sql.Stmt
	tableCreated bool
}

// Init creates table if needed
func (c *Calendar) createTableOnce() {
	if c.tableCreated {
		return
	}
	_, err := c.DB.Exec(`
    create table Calendar (
      Date datetime not null
    )
  `)
	if err == nil {
		c.tableCreated = true
		glog.Info("Calendar table created")
	} else {
		glog.Warning(err)
	}
}

func (c *Calendar) Add(dates ...date) error {
	c.createTableOnce()

	if c.add == nil {
		var err error
		c.add, err = c.DB.Prepare(`insert into Calendar (Date) values ($1)`)
		if err != nil {
			panic(err)
		}
	}
	for _, d := range dates {
		_, err := c.add.Exec(d)
		if err != nil {
			glog.Error(err)
			return err
		}
	}
	return nil
}

func (c *Calendar) Available(start, stop date) (bool, error) {
	list, err := c.List()
	if err != nil {
		return false, err
	}

	var include = func(l []date, current date) bool {
		for _, d := range l {
			if d == current {
				return true
			}
		}
		return false
	}

	daysApart := start.DaysApart(stop)
	glog.Warningln("daysApart", daysApart)
	for i := 0; i <= daysApart; i++ {
		current := start.Add(i)
		glog.Warningln("current", current)
		if !include(list, current) {
			return false, unavailable
		}
	}

	return true, nil
}

func (c *Calendar) List() ([]date, error) {
	c.createTableOnce()

	if c.list == nil {
		var err error
		c.list, err = c.DB.Prepare(`select Date from Calendar`)
		if err != nil {
			panic(err)
		}
	}

	rows, err := c.list.Query()
	if err != nil {
		glog.Error(err)
		return []date{}, err
	}
	defer rows.Close()

	var list []date
	for rows.Next() {
		var d date
		err := rows.Scan(&d)
		if err != nil {
			glog.Error(err)
			return []date{}, err
		}
		list = append(list, d)
	}

	return list, nil
}

func (c *Calendar) Remove(dates ...date) error {
	c.createTableOnce()

	if c.remove == nil {
		var err error
		c.remove, err = c.DB.Prepare(`delete from Calendar where Date=$1`)
		if err != nil {
			panic(err)
		}
	}
	for _, d := range dates {
		_, err := c.remove.Exec(d)
		if err != nil {
			glog.Error(err)
			return err
		}
	}
	return nil
}
