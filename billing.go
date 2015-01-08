package booking

import "time"

type charger interface {
	Charge(creditCard, amount) error
}

type creditCard struct {
	Name            string
	Number          string
	ExpirationMonth int
	ExpirationYear  int
}

type paymentProcessor interface {
	Authorize(creditCard) (paymentToken, error)
}

type paymentToken struct {
	token string
}

func (pt paymentToken) Capture(amount) error {
	return nil
}

const minAmountCents = 100    // $1.00
const maxAmountCents = 100000 // $1,000.00

type amount struct {
	Cents int // positive or negative
}

type entry struct {
	Amount amount
	Date   time.Time
}

type ledger map[guestId][]entry

func (l ledger) Credit(guestId, amount) error { return nil }
func (l ledger) Debit(guestId, amount) error  { return nil }

type rateCode int

const (
	rateWithBunny rateCode = iota
	rateWithoutBunny
)

func (c rateCode) Amount() amount {
	switch c {
	case rateWithoutBunny:
		return amount{25000}
	case rateWithBunny:
		return amount{10000}
	default:
		panic("unknown rate code")
	}
}
