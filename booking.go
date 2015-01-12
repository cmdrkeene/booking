package booking

import (
	"errors"
	"fmt"
	"log"
	"sync/atomic"
	"time"
)

type bookingId uint32

func (id bookingId) String() string {
	return fmt.Sprintf("booking:%d ", id)
}

// Workflow to orchestrate a guest paying for a reservation
type booking struct {
	created time.Time
	dates   dateRange
	guestId guestId
	id      bookingId
	rate    rateCode
	state   bookingState

	// dependencies
	charger   charger
	registrar registrar
	reserver  reserver
	store     bookingStore
}

type bookingStore interface {
	NewId() bookingId
	Save(*booking) error
}

type bookingMemoryStore struct {
	lastId  uint32
	records map[bookingId]*booking
}

func newBookingMemoryStore() *bookingMemoryStore {
	return &bookingMemoryStore{
		lastId:  0,
		records: make(map[bookingId]*booking),
	}
}

func (s *bookingMemoryStore) NewId() bookingId {
	return bookingId(atomic.AddUint32(&s.lastId, 1))
}

func (s *bookingMemoryStore) Save(b *booking) error {
	s.records[b.id] = b
	return nil
}

func newBooking(
	charger charger,
	registrar registrar,
	reserver reserver,
	store bookingStore,
) (*booking, error) {
	// workflow
	b := &booking{}
	b.created = time.Now()
	b.id = store.NewId()
	b.state = bookingScheduling

	// dependencies
	b.charger = charger
	b.registrar = registrar
	b.reserver = reserver
	b.store = store

	err := store.Save(b)
	if err != nil {
		return nil, err
	}
	log.Print(b.id, "created")
	return b, nil
}

func (b *booking) Schedule(dates dateRange, rate rateCode) error {
	if b.state != bookingScheduling {
		return bookingStateError
	}

	if !b.reserver.IsAvailable(dates, rate) {
		return unavailable
	}

	b.dates = dates
	b.rate = rate
	b.state = bookingRegistering

	err := b.store.Save(b)
	if err != nil {
		return err
	}

	log.Print(b.id, "scheduled")
	return nil
}

func (b *booking) Register(name, email string) error {
	if b.state != bookingRegistering {
		return bookingStateError
	}

	guestId, err := b.registrar.Register(name, email)
	if err != nil {
		return err
	}

	b.guestId = guestId
	b.state = bookingPaying

	err = b.store.Save(b)
	if err != nil {
		return err
	}
	log.Print(b.id, "registered")
	return nil
}

func (b *booking) Pay(card creditCard) error {
	if b.state != bookingPaying {
		return bookingStateError
	}

	amount := b.rate.Amount()
	err := b.charger.Charge(card, amount)
	if err != nil {
		return err
	}

	b.state = bookingReserving
	err = b.store.Save(b)
	if err != nil {
		return err
	}
	log.Print(b.id, "paid")
	return nil
}

func (b *booking) Reserve() error {
	if b.state != bookingReserving {
		return bookingStateError
	}

	err := b.reserver.Reserve(b.dates, b.rate, b.guestId)
	if err != nil {
		return err
	}

	b.state = bookingComplete
	err = b.store.Save(b)
	if err != nil {
		return err
	}

	log.Print(b.id, "complete")
	return nil
}

type bookingState int

// scheduling -> registering -> paying -> reserving -> complete
const (
	bookingScheduling bookingState = iota
	bookingRegistering
	bookingPaying
	bookingReserving
	bookingComplete
)

var bookingStateError = errors.New("bad booking state")

func (s bookingState) String() string {
	switch s {
	case bookingScheduling:
		return "scheduling"
	case bookingRegistering:
		return "registering"
	case bookingPaying:
		return "paying"
	case bookingReserving:
		return "reserving"
	case bookingComplete:
		return "complete"
	default:
		panic("unknown booking state")
	}
}
