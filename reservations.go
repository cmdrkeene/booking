package booking

import (
	"errors"
	"strings"
	"sync/atomic"
	"time"
)

// Reserves dates for a guest for a price
type Reserver interface {
	Reserve(guestId, dateRange, rateCode) error
}

type reservationId string

type reservation struct {
	Id      reservationId
	Created time.Time
	GuestId guestId
	Dates   dateRange
	Rate    rateCode
}

// if event exists and does not have a reservationId, it is "available"
type event struct {
	reservationId reservationId
}

func (e event) Available() bool {
	return !e.Reserved()
}

func (e event) Reserved() bool {
	return len(e.reservationId) > 0
}

type calendar struct {
	events                map[time.Time]event
	reservations          []reservation
	reservationPrimaryKey uint32
}

func newCalendar() calendar {
	return calendar{
		events: make(map[time.Time]event),
	}
}

func (c calendar) String() string {
	var lines []string
	lines = append(lines, "\n== Calendar ==")
	for t, event := range c.events {
		l := t.Format(dayFormat)
		if event.Reserved() {
			l = l + " (Reserved)"
		}
		lines = append(lines, l)
	}
	return strings.Join(lines, "\n")
}

func (c calendar) SetAvailable(r dateRange) {
	for _, t := range r.EachDay() {
		c.events[t] = event{}
	}
}

var unavailable = errors.New("dates unavailable")

func (c calendar) Reserve(gid guestId, dr dateRange, rc rateCode) error {
	// are all days available?
	for _, t := range dr.EachDay() {
		event, ok := c.events[t]
		if !ok {
			return unavailable
		}
		if !event.Available() {
			return unavailable
		}
	}

	// new reservation
	id := c.newReservation(gid, dr, rc)

	// mark on calendar
	for _, t := range dr.EachDay() {
		c.events[t] = event{reservationId: id}
	}
	return nil
}

func (c calendar) newReservation(g guestId, dr dateRange, rc rateCode) reservationId {
	res := reservation{
		Id:      c.newReservationId(),
		Created: time.Now(),
		Dates:   dr,
		GuestId: g,
		Rate:    rc,
	}
	c.reservations = append(c.reservations, res)
	return res.Id
}

func (c calendar) newReservationId() reservationId {
	id := atomic.AddUint32(&c.reservationPrimaryKey, 1)
	return reservationId(id)
}

const day = 24 * time.Hour
const dayFormat = "January 2, 2006"

type dateRange struct {
	days  int // number of days from start, must be > 0
	start time.Time
}

func newDateRange(t time.Time, days int) dateRange {
	if days == 0 {
		panic("minimum days is 1")
	}
	return dateRange{start: t, days: days}
}

func (r dateRange) EachDay() []time.Time {
	var days []time.Time
	for i := 0; i < r.days; i++ {
		delta := time.Duration(i) * day
		days = append(days, r.start.Add(delta))
	}
	return days
}

func (r dateRange) String() string {
	t1 := r.start.Format(dayFormat)
	t2 := r.start.Add(time.Duration(r.days) * day).Format(dayFormat)
	return t1 + " to " + t2
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
