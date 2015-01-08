package booking

import (
	"testing"
	"time"
)

// scheduling -> registering -> paying -> complete
func TestBookingStateMachine(t *testing.T) {
	booking := newBooking(guestId(123))
	if booking.state != bookingScheduling {
		t.Error("want", bookingScheduling)
		t.Error("got ", booking.state)
	}

	// Schedule
	dates := newDateRange(time.Now(), 1)
	err := booking.Schedule(
		dates,
		rateWithBunny,
		testReserver{err: nil},
	)
	if err != nil {
		t.Error(err)
	}
	if booking.state != bookingRegistering {
		t.Error("want", bookingRegistering)
		t.Error("got ", booking.state)
	}

	// Register
	err = booking.Register(
		"Brandon",
		"brandon@example.com",
		testRegistrar{err: nil, guestId: guestId(999)},
	)
	if err != nil {
		t.Error(err)
	}
	if booking.state != bookingPaying {
		t.Error("want", bookingPaying)
		t.Error("got ", booking.state)
	}

	// Pay
	err = booking.Pay(creditCard{}, testCharger{})
	if err != nil {
		t.Error(err)
	}
	if booking.state != bookingReserving {
		t.Error("want", bookingReserving)
		t.Error("got ", booking.state)
	}

	// Reserve
	err = booking.Reserve(testReserver{})
	if err != nil {
		t.Error(err)
	}
	if booking.state != bookingComplete {
		t.Error("want", bookingComplete)
		t.Error("got ", booking.state)
	}

	// Complete - check final booking data
	if booking.dates != dates {
		t.Error("want", dates)
		t.Error("got ", booking.dates)
	}

	if booking.rate != rateWithBunny {
		t.Error("want", rateWithBunny)
		t.Error("got ", booking.rate)
	}

	if booking.guestId != guestId(999) {
		t.Error("want", guestId(999))
		t.Error("got ", booking.guestId)
	}
}

type testCharger struct{ err error }

func (t testCharger) Charge(creditCard, amount) error {
	return t.err
}

type testRegistrar struct {
	err     error
	guestId guestId
}

func (t testRegistrar) Register(name, email string) (guestId, error) {
	if t.err == nil {
		return t.guestId, nil
	}
	return 0, t.err
}

type testReserver struct{ err error }

func (t testReserver) IsAvailable(dateRange, rateCode) bool {
	return t.err == nil
}

func (t testReserver) Reserve(dateRange, rateCode, guestId) error {
	return t.err
}
