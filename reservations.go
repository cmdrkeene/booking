package booking

import (
	"strings"
	"time"
)

// Manages room inventory for dates
type Reservations interface {
	Reserve(guestId, dateRange, rate) error
}

type reservationId string

// Record of someone reserving dates for a fee
type reservation struct {
	Id      reservationId
	Created time.Time
	GuestId guestId
	Dates   dateRange
	Rate    rate
}

type rate struct {
	Amount amount
	Name   string
}

var rateComplimentary = rate{Amount: amount{0}, Name: "Complimentary"}
var rateWithBunny = rate{Amount: amount{10000}, Name: "With Bunny"}
var rateWithoutBunny = rate{Amount: amount{25000}, Name: "Without Bunny"}

type reservationService struct{}

func (rs reservationService) Reserve(guestId, dateRange, rate) (reservationId, error) {
	return reservationId{}, nil
}

type calendar map[time.Time]availability

type availability struct {
	Booked bool
}

func (c calendar) String() string {
	var lines []string
	lines = append(lines, "\n== Calendar ==")
	for t, a := range c {
		l := t.Format(dayFormat)
		if a.Booked {
			l = l + " (Booked)"
		}
		lines = append(lines, l)
	}
	return strings.Join(lines, "\n")
}

func (c calendar) SetAvailable(r dateRange) {
	for _, t := range r.Days() {
		c[t] = availability{Booked: false}
	}
}

func (c calendar) SetBooked(r dateRange) bool {
	// check if all available and not booked
	for _, t := range r.Days() {
		a, ok := c[t]
		if !ok {
			return false // unavailable
		}
		if a.Booked {
			return false // booked
		}
	}

	// mark it
	for _, t := range r.Days() {
		c[t] = availability{Booked: true}
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
