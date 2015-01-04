/* Actions

As an admin,

- mark dates as available, unavailable
- view any offers
- accept an offer (creates an availability, booking, marks offer as accepted/hidden)
- decline an offer (deletes offer, sends decline email)

As a guest

- submit an offer
- make a booking (pay for it too)

*/
package booking

import "time"

const RateWithBunny = 150
const RateWithoutBunny = 250

type dateRange struct {
	Start    time.Time
	Duration time.Duration
}

// order of dates doesn't matter. Start is earlier time.
func newDateRange(t1, t2 time.Time) dateRange {
	if t1.Before(t2) {
		return dateRange{Start: t1, Duration: t2.Sub(t1)}
	} else {
		return dateRange{Start: t2, Duration: t1.Sub(t2)}
	}
}

type bookingService struct{}

func (s bookingService) Book(guest, dateRange) (booking, error) {
	return booking{}, nil
}
func (s bookingService) Purchase(payment, booking) error { return nil }
func (s bookingService) Offer(swap) error                { return nil }

type calendar struct{}

func (c calendar) IsAvailable(dateRange) bool    { return false }
func (c calendar) IsBusy(dateRange) bool         { return false }
func (c calendar) MarkAvailable(dateRange) error { return nil }
func (c calendar) MarkBusy(dateRange) error      { return nil }

type guestService struct {
	guests []guest
}

func (s guestService) Register(name, email string) (guest, error) {
	return guest{}, nil
}

// When we offer the apartment
type availability struct {
	Is    id
	Dates dateRange
	Title string // "Labor Day Weekend"
}

type id string

// When someone has booked an availabilt
type booking struct {
	Id             id
	AvailabilityId id
	Bunny          bool // will they watch the bunny?
	Payment        payment
	Guest          guest
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

type payment struct {
	Cents          int64
	Date           time.Time
	ProcessorToken string // from stripe, paypal etc. should have type, last 4, etc
}

// The rando sleeping/having-sex in your bed
type guest struct {
	Name  string // required
	Email string // required
}
