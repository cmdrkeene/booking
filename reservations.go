package booking

import (
	"errors"
	"log"
	"sync/atomic"
	"time"
)

var unavailable = errors.New("unavailable")

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
	List() ([]reservation, error)
}

type reserver interface {
	Available() []time.Time
	IsAvailable(dateRange) bool
	Reserve(dateRange, rateCode, guestId) error
}

type reservationManager struct {
	availability availabilityStore
	reservations reservationStore
}

func newReservationManager(a availabilityStore, r reservationStore) reservationManager {
	return reservationManager{a, r}
}

func (m reservationManager) Reserve(dr dateRange, rc rateCode, id guestId) error {
	// check available
	available, err := m.availability.List()
	if err != nil {
		return err
	}
	log.Print("availability", available)
	if !dr.Coincident(available) {
		return unavailable
	}

	// check reserved
	reserved, err := m.reserved()
	log.Print("reserved", reserved)
	if err != nil {
		return err
	}
	if dr.Coincident(reserved) {
		return unavailable
	}

	// save
	record := &reservation{Dates: dr, GuestId: id, Rate: rc}
	err = m.reservations.Save(record)
	if err != nil {
		return err
	}

	return nil
}

func (m reservationManager) reserved() ([]time.Time, error) {
	list, err := m.reservations.List()
	if err != nil {
		return []time.Time{}, err
	}

	var times []time.Time
	for _, r := range list {
		for _, t := range r.Dates.EachDay() {
			times = append(times, t)
		}
	}
	return times, nil
}

type reservationMemoryStore struct {
	lastId  uint32
	records map[reservationId]*reservation
}

func newReservationMemoryStore() reservationMemoryStore {
	return reservationMemoryStore{
		lastId:  0,
		records: make(map[reservationId]*reservation),
	}
}

func (s reservationMemoryStore) newId() reservationId {
	return reservationId(atomic.AddUint32(&s.lastId, 1))
}

func (s reservationMemoryStore) List() ([]reservation, error) {
	var list []reservation
	for _, rec := range s.records {
		list = append(list, *rec)
	}
	return list, nil
}

func (s reservationMemoryStore) Save(record *reservation) error {
	if record.Id == 0 {
		record.Id = s.newId()
	}
	if record.Created.IsZero() {
		record.Created = time.Now()
	}
	s.records[record.Id] = record
	return nil
}
