/* Actions

As an admin,

- mark dates as available, unavailable
- view any offers
- accept an offer (creates an availability, booking, marks offer as accepted/hidden)
- decline an offer (deletes offer, sends decline email)

As a guest

- submit an offer
- make a booking (pay for it too)

## Departments

Guestbook
\_Guest - a record of a person staying at the hotel

Booking - manages inventory (room(s) on dates)
\_Calendar - a master record of available/confirmed dates
\_Reservation - a record of guest's payment for dates

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

const rateWithBunny = 150
const rateWithoutBunny = 250
const day = 24 * time.Hour

type sessionId string
type sessionState int

var sessionStateError = errors.New("invalid session state")

const (
	sessionStarted sessionState = iota
	sessionNeedsPayment
	sessionComplete
)

// state machine to orchestrates guest, payment, and booking
type session struct {
	created time.Time
	dates   dateRange
	guestId guestId
	id      sessionId
	state   sessionState
}

func newSession(guestId) *session {
	s := &session{
		id:      sessionId("new"),
		created: time.Now(),
		guestId: guestId,
		state:   sessionStarted,
	}
	log.Print("session creates", s)
	return s
}

func (s *session) ChooseDates(r dateRange) error {
	if s.state != sessionStarted {
		return errors.New("invalid session state")
	}

	// check availability

	// mark as needsPayment
	s.state = sessionNeedsPayment
	log.Print("chose dates", r)
	return nil
}

func (s *session) Pay(creditCard) error {
	if s.state != sessionNeedsPayment {
		return errors.New("session must be needs payment")
	}

	// tell billing to capture payment
	// tell booking to confirm reservation
}

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

const dayFormat = "January 2, 2006"

func (r dateRange) String() string {
	t1 := r.Start.Format(dayFormat)
	t2 := r.Start.Add(time.Duration(r.NumDays) * day).Format(dayFormat)
	return t1 + " to " + t2
}

type bookingService struct{}

func (s bookingService) Book(guest, dateRange) error {
	return nil
}

func (s bookingService) Offer(swap) error { return nil }

type booking struct {
	Bunny bool // will they watch the bunny?
	Guest guest
}

// TODO mark accepted or delete / log / audit
type swap struct {
	Address     string
	Attachments []byte
	Bunny       bool
	Dates       dateRange
	Description string // roomate apartment, studio, roomates, images
	Guest       guest
}
