package booking

import (
	"database/sql/driver"
	"errors"
	"fmt"
	"time"
)

// Simple date
type date struct {
	t time.Time
}

const iso8601Date = "2006-01-02"

func (d date) String() string {
	return d.t.Format(iso8601Date)
}

func (d *date) Scan(src interface{}) error {
	t, ok := src.(time.Time)
	if !ok {
		err := errors.New(
			fmt.Sprintf("can't scan date from db: %#v", src),
		)
		return err
	}
	d.t = t.UTC()
	return nil
}

func (d date) Value() (driver.Value, error) {
	return driver.Value(d.t), nil
}

func newDate(year, month, day int) date {
	return date{
		t: time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.UTC),
	}
}
