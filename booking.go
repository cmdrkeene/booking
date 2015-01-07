/* Actions

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

const rateWithBunny = 150
const rateWithoutBunny = 250
const day = 24 * time.Hour

// state machine to orchestrates a guest paying for a reservation
type bookingWorkflow struct {
	created time.Time
	dates   dateRange
	guestId guestId
	state   bookingWorkflowState
}

type bookingWorkflowState int

const (
	bookingStarted bookingWorkflowState = iota
	bookingNeedsPayment
	bookingFinished
)

func (w *bookingWorkflow) IsStarted() bool      { return w.state == bookingStarted }
func (w *bookingWorkflow) IsNeedsPayment() bool { return w.state == bookingNeedsPayment }
func (w *bookingWorkflow) IsFinished() bool     { return w.state == bookingFinished }

func startBookingWorkflow(g guestId) *bookingWorkflow {
	w := &bookingWorkflow{
		created: time.Now(),
		guestId: g,
		state:   bookingStarted,
	}
	log.Print("booking workflow started", w)
	return w
}

func (w *bookingWorkflow) ChooseDates(r dateRange) error {
	if !w.IsStarted() {
		return errors.New("invalid session state")
	}

	// check availability

	// mark as needsPayment
	w.state = bookingNeedsPayment
	log.Print("chose dates", r)
	return nil
}

func (w *bookingWorkflow) Pay(creditCard) error {
	if !w.IsNeedsPayment() {
		return errors.New("invalid session state")
	}

	// tell billing to capture payment
	// tell booking to confirm reservation
	return nil
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
