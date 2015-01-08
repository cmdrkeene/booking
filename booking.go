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
	rate    rateCode
	state   bookingState
}

func newBooking(g guestId) *booking {
	w := &booking{
		created: time.Now(),
		guestId: g,
		state:   bookingScheduling,
	}
	return w
}

func (b *booking) Schedule(dates dateRange, rate rateCode, res reserver) error {
	if b.state != bookingScheduling {
		return bookingStateError
	}

	if !res.IsAvailable(dates, rate) {
		return unavailable
	}

	b.dates = dates
	b.rate = rate
	b.state = bookingRegistering
	return nil
}

func (b *booking) Register(name, email string, reg registrar) error {
	if b.state != bookingRegistering {
		return bookingStateError
	}

	guestId, err := reg.Register(name, email)
	if err != nil {
		return err
	}

	b.guestId = guestId
	b.state = bookingPaying
	return nil
}

func (b *booking) Pay(card creditCard, chg charger) error {
	if b.state != bookingPaying {
		return bookingStateError
	}

	amount := b.rate.Amount()
	err := chg.Charge(card, amount)
	if err != nil {
		return err
	}

	b.state = bookingReserving
	return nil
}

func (b *booking) Reserve(res reserver) error {
	if b.state != bookingReserving {
		return bookingStateError
	}

	err := res.Reserve(b.dates, b.rate, b.guestId)
	if err != nil {
		return err
	}

	b.state = bookingComplete

	return nil
}

type bookingState int

// scheduling -> registering -> paying -> reserving - complete
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
