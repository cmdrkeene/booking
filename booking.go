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
	Duration time.Time
}

type bookingService struct{}

func (s bookingService) Book(guest, dateRange) (booking, error) {}
func (s bookingService) Purchase(payment, booking) error        {}
func (s bookingService) Offer(swap) error                       {}

type calendar struct{}

func (c calendar) IsAvailable(dateRange) bool
func (c calendar) IsBusy(dateRange) bool
func (c calendar) MarkAvailable(dateRange) error
func (c calendar) MarkBusy(dateRange) error

type guestService struct {
	guests []guest
}

func (s guestService) Register(name, email string) (guest, error) {}

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
