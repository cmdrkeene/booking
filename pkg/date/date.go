// Package Date makes working with Dates easier than time.Time
// Date is stored in SQL as a time.Time
package date

import (
	"database/sql/driver"
	"errors"
	"fmt"

	"time"
)

// Duration helper
const Day = 24 * time.Hour

// Formats
const (
	Pretty  = "January 2, 2006"
	ISO8601 = "2006-01-02"
)

// Simpler Date object than time.Time
type Date struct {
	t time.Time
}

func New(year, month, day int) Date {
	return Date{
		t: time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.UTC),
	}
}

var ParseError = errors.New("can't parse Date")

func Parse(src interface{}) (Date, error) {
	switch t := src.(type) {
	default:
		return Date{}, ParseError
	case string:
		return parseString(t)
	case time.Time:
		return parseTime(t)
	}
}

func parseString(s string) (Date, error) {
	if s == "" {
		return Date{}, ParseError
	}

	var t time.Time
	var err error

	dateFormats := []string{
		"01/02/2006",
		"01/02/06",
		"1/2/2006",
		"1/2/06",
	}
	for _, format := range dateFormats {
		t, err = time.Parse(format, s)
		if err == nil {
			return parseTime(t)
		}
	}

	t, err = time.Parse(ISO8601, s)
	if err == nil {
		return parseTime(t)
	}

	return Date{}, ParseError
}

func parseTime(t time.Time) (Date, error) {
	return New(t.Year(), int(t.Month()), t.Day()), nil
}

func (d Date) Add(n int) Date {
	return Date{
		d.t.Add(time.Duration(n) * Day),
	}
}

func (d Date) After(u Date) bool {
	return d.t.After(u.t)
}

func (d Date) DaysApart(u Date) int {
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

func (d Date) Format(layout string) string {
	return d.t.Format(layout)
}

func (d *Date) Scan(src interface{}) error {
	t, ok := src.(time.Time)
	if !ok {
		err := errors.New(
			fmt.Sprintf("can't scan Date from sql: %#v", src),
		)
		return err
	}
	d.t = t.UTC()
	return nil
}

func (d Date) String() string {
	return d.Format(ISO8601)
}

func (d Date) Value() (driver.Value, error) {
	return driver.Value(d.t), nil
}
