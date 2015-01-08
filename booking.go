package booking

import (
	"errors"
	"time"
)

// Workflow to orchestrate a guest paying for a reservation
type booking struct {
	created time.Time
	dates   dateRange
	guestId guestId
	state   bookingState
}

// Scheduling -> Registering -> Paying -> Complete
type bookingState int

const (
	bookingScheduling bookingState = iota
	bookingRegistering
	bookingPaying
	bookingComplete
)

// invalid state transition
var bookingStateError = errors.New("bad booking state")

func newBooking(g guestId) *booking {
	w := &booking{
		created: time.Now(),
		guestId: g,
		state:   bookingScheduling,
	}
	return w
}

func (b *booking) Schedule(dates dateRange, r Reserver) error {
	if b.state != bookingScheduling {
		return bookingStateError
	}
	// check availability
	b.state = bookingRegistering
	return nil
}

func (b *booking) Register(name, email string) error {
	if b.state != bookingScheduling {
		return bookingStateError
	}

	b.state = bookingPaying
	return nil
}

func (b *booking) Pay(c creditCard, billing Billing) error {
	if b.state != bookingPaying {
		return bookingStateError
	}

	// tell billing to capture payment
	// tell booking to confirm reservation
	b.state = bookingComplete
	return nil
}
