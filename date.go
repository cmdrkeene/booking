package booking

import (
	"database/sql/driver"
	"errors"
	"fmt"
	"time"

	"github.com/golang/glog"
)

const day = 24 * time.Hour

// Make working with dates easier than time.Time
type date struct {
	t time.Time
}

func newDate(year, month, day int) date {
	return date{
		t: time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.UTC),
	}
}

var invalidDate = errors.New("invalid date")

func newDateFromString(s string) (date, error) {
	if s == "" {
		return date{}, invalidDate
	}

	t, err := time.Parse(layoutDateISO8601, s)
	if err != nil {
		return date{}, invalidDate
	}

	return newDate(t.Year(), int(t.Month()), t.Day()), nil
}

func (d date) Add(n int) date {
	return date{
		d.t.Add(time.Duration(n) * day),
	}
}

func (d date) After(u date) bool {
	return d.t.After(u.t)
}

func (d date) DaysApart(u date) int {
	duration := d.t.Sub(u.t)

	// Didn't seem to be a stdlib way to do this
	var abs = func(x int) int {
		if x < 0 {
			return -x
		}
		return x
	}

	return abs(int(duration.Hours() / 24))
}

const (
	layoutDatePretty  = "January 2, 2006"
	layoutDateISO8601 = "2006-01-02"
)

func (d date) Format(layout string) string {
	return d.t.Format(layout)
}

func (d *date) Scan(src interface{}) error {
	t, ok := src.(time.Time)
	if !ok {
		err := errors.New(
			fmt.Sprintf("can't scan date from db: %#v", src),
		)
		glog.Error(err)
		return err
	}
	d.t = t.UTC()
	return nil
}

func (d date) String() string {
	return d.Format(layoutDateISO8601)
}

func (d date) Value() (driver.Value, error) {
	return driver.Value(d.t), nil
}
