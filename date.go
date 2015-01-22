package booking

import (
	"database/sql/driver"
	"errors"
	"fmt"
	"time"

	"github.com/golang/glog"
)

const (
	day         = 24 * time.Hour
	iso8601Date = "2006-01-02"
)

// Make working with dates easier than time.Time
type date struct {
	t time.Time
}

func newDate(year, month, day int) date {
	return date{
		t: time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.UTC),
	}
}

func newDateFromString(s string) (date, error) {
	t, err := time.Parse(iso8601Date, s)
	if err != nil {
		return date{}, err
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

func (d date) String() string {
	return d.t.Format(iso8601Date)
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

func (d date) Value() (driver.Value, error) {
	return driver.Value(d.t), nil
}
