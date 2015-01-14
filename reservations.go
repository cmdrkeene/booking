package booking

import (
	"errors"
	"strings"
	"sync/atomic"
	"time"
)

type reservationId uint32

type reservation struct {
	Id      reservationId
	Created time.Time
	GuestId guestId
	Dates   dateRange
	Rate    rateCode
}

type reservationStore interface {
	Save(*reservation) error
}

type reserver interface {
	Available() []time.Time
	IsAvailable(dateRange) bool
	Reserve(dateRange, rateCode, guestId) error
}

type reservationMemoryStore struct {
	lastId  uint32
	records map[reservationId]*reservation
}

func newReservationMemoryStore() *reservationMemoryStore {
	return &reservationMemoryStore{
		lastId:  0,
		records: make(map[reservationId]*reservation),
	}
}

func (s *reservationMemoryStore) newId() reservationId {
	return reservationId(atomic.AddUint32(&s.lastId, 1))
}

func (s *reservationMemoryStore) Save(record *reservation) error {
	if record.Id == 0 {
		record.Id = s.newId()
	}
	s.records[record.Id] = record
	return nil
}

type event struct {
	reservationId reservationId
}

func (e event) IsReserved() bool {
	return e.reservationId > 0
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
		l := t.Format(pretty)
		if event.IsReserved() {
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

func (c calendar) Available() []time.Time {
	var days []time.Time
	for t, e := range c.events {
		if !e.IsReserved() {
			days = append(days, t)
		}
	}
	return days
}

func (c calendar) IsAvailable(dr dateRange) bool {
	// are all days available?
	for _, t := range dr.EachDay() {
		event, ok := c.events[t]
		if !ok {
			return false
		}
		if event.IsReserved() {
			return false
		}
	}
	return true
}

func (c calendar) Reserve(dr dateRange, rc rateCode, gid guestId) error {
	if !c.IsAvailable(dr) {
		return unavailable
	}

	id := c.newReservation(gid, dr, rc)
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
