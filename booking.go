/* Actions

TODO define top level interfaces for services
TODO test reservation service (crud for reservation records)
TODO test billing service with dummy processor
TODO test workflow state transitions, errors from services
TODO basic controller/pages for workflow
TODO workflow to offer aparment swap
TODO admin workflow

As an admin,

- mark dates as available, unavailable
- view any offers
- accept an offer (creates an availability, booking, marks offer as accepted/hidden)
- decline an offer (deletes offer, sends decline email)

As a guest

- submit an offer
- make a booking (pay for it too)

## Departments/Servives

BookingWorkflow - a coordinator that requires payment for bookings

Guestbook
\_Guest - a record of a person staying at the hotel

Reservation - manages inventory (room(s) on dates)
\_Calendar - a master record of available/confirmed dates
\_Booking - a record of guest's payment for dates

Billing - manages bookeeping, credit cards
\_CreditCard - unsaved raw card information
\_Ledger - credits and debits for a guest's account
\_Processor - entity that manages credit cards
\_PaymentToken - persisted credit card proxy

## Dependencies


// what is the thing that orchestrates it
bookingApp := newBookingApp(
	newBookingService(),
	newBillingService(),
	newGuestService(),
)

bookingApp.LogIn()
bookingApp.LogOut()
bookingApp.ChooseDates()
bookingApp.Pay()

*/
package booking

import (
	"errors"
	"log"
	"time"
)

// Workflow to orchestrate a guest paying for a reservation
type booking struct {
	created time.Time
	dates   dateRange
	guestId guestId
	state   bookingState
}

type bookingState int

const (
	bookingStarted bookingState = iota
	bookingNeedsPayment
	bookingFinished
)

func (b *booking) IsStarted() bool      { return w.state == bookingStarted }
func (b *booking) IsNeedsPayment() bool { return w.state == bookingNeedsPayment }
func (b *booking) IsFinished() bool     { return w.state == bookingFinished }

func startBooking(g guestId) *bookingWorkflow {
	w := &booking{
		created: time.Now(),
		guestId: g,
		state:   bookingStarted,
	}
	log.Print("booking workflow started", w)
	return w
}

func (b *booking) ChooseDates(r dateRange) error {
	if !b.IsStarted() {
		return errors.New("invalid session state")
	}

	// check availability

	// mark as needsPayment
	b.state = bookingNeedsPayment
	log.Print("chose dates", r)
	return nil
}

func (b *booking) Pay(creditCard) error {
	if !w.IsNeedsPayment() {
		return errors.New("invalid session state")
	}

	// tell billing to capture payment
	// tell booking to confirm reservation
	return nil
}
