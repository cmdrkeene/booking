package booking

import (
	"strings"
	"time"
)

// Manages room inventory for dates
type Reservations interface {
	Reserve(guestId, dateRange, rateCode) error
}

type reservationId string

// Record of someone reserving dates for a fee
type reservation struct {
	Id      reservationId
	Created time.Time
	GuestId guestId
	Dates   dateRange
	Rate    rateCode
}

type rateCode int

const (
	rateWithBunny rateCode = iota
	rateWithoutBunny
)

func rateAmount(c rateCode) amount {
	switch c {
	case rateWithoutBunny:
		return amount{25000}
	case rateWithBunny:
		return amount{10000}
	default:
		panic("unknown rate code")
	}
}

// sparse map of available or booked dates
type calendar map[time.Time]event

// if event is on calendar, it is assumed to be available
type event struct {
	reservationId reservationId
}

func (e event) Available() bool {
	return !e.Reserved()
}

func (e event) Reserved() bool {
	return len(e.reservationId) > 0
}

func (c calendar) String() string {
	var lines []string
	lines = append(lines, "\n== Calendar ==")
	for t, event := range c {
		l := t.Format(dayFormat)
		if event.Reserved() {
			l = l + " (Reserved)"
		}
		lines = append(lines, l)
	}
	return strings.Join(lines, "\n")
}

func (c calendar) SetAvailable(r dateRange) {
	for _, t := range r.Days() {
		c[t] = event{}
	}
}

func (c calendar) Reserve(dr dateRange, ri reservationId) bool {
	// check if all available and not booked
	for _, t := range dr.Days() {
		event, ok := c[t]
		if !ok {
			return false // not available
		}
		if event.Reserved() {
			return false
		}
	}

	// mark it
	for _, t := range dr.Days() {
		c[t] = event{reservationId: ri}
	}
	return true
}

const day = 24 * time.Hour
const dayFormat = "January 2, 2006"

type dateRange struct {
	NumDays int
	Start   time.Time
}

func newDateRange(start time.Time, numDays int) dateRange {
	return dateRange{Start: start, NumDays: numDays}
}

func (r dateRange) Days() []time.Time {
	var days []time.Time
	for i := 0; i < r.NumDays; i++ {
		delta := time.Duration(i) * day
		days = append(days, r.Start.Add(delta))
	}
	return days
}

func (r dateRange) String() string {
	t1 := r.Start.Format(dayFormat)
	t2 := r.Start.Add(time.Duration(r.NumDays) * day).Format(dayFormat)
	return t1 + " to " + t2
}
