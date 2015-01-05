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

import (
	"strings"
	"time"
)

const rateWithBunny = 150
const rateWithoutBunny = 250
const day = 24 * time.Hour

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

func (s bookingService) Book(guest, dateRange) (booking, error) {
	return booking{}, nil
}
func (s bookingService) Purchase(payment, booking) error { return nil }
func (s bookingService) Offer(swap) error                { return nil }

type availability struct {
	Booked bool
}

type calendar map[time.Time]availability

func (c calendar) String() string {
	var lines []string
	lines = append(lines, "\n== Calendar ==")
	for t, a := range c {
		l := t.Format(dayFormat)
		if a.Booked {
			l = l + " (Booked)"
		}
		lines = append(lines, l)
	}
	return strings.Join(lines, "\n")
}

func (c calendar) SetAvailable(r dateRange) {
	for _, t := range r.Days() {
		c[t] = availability{Booked: false}
	}
}

func (c calendar) SetBooked(r dateRange) bool {
	// check if all available and not booked
	for _, t := range r.Days() {
		a, ok := c[t]
		if !ok {
			return false // unavailable
		}
		if a.Booked {
			return false // booked
		}
	}

	// mark it
	for _, t := range r.Days() {
		c[t] = availability{Booked: true}
	}
	return true
}

func (c calendar) Unbook(r dateRange) {}

type guestService struct {
	guests []guest
}

func (s guestService) Register(name, email string) (guest, error) {
	return guest{}, nil
}

type id string

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
