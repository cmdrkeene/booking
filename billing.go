package booking

import (
	"fmt"
	"log"
	"time"

	"github.com/stripe/stripe-go"
	"github.com/stripe/stripe-go/charge"
)

type charger interface {
	Charge(creditCard, amount) error
}

type fakeCharger struct{}

func (f fakeCharger) Charge(c creditCard, a amount) error {
	log.Println("Charging", c, a.InDollars())
	return nil
}

type stripeCharger struct{}

const stripeTestKey = "sk_test_BQokikJOvBiI2HlWgH4olfQ2"

func newStripeCharger(string) stripeCharger {
	stripe.Key = stripeTestKey
	return stripeCharger{}
}

func (s stripeCharger) Charge(c creditCard, a amount) error {
	log.Println("Charging", c, a.InDollars())
	params := &stripe.ChargeParams{
		Amount:   uint64(a.Cents),
		Currency: "usd",
		Card: &stripe.CardParams{
			Number: c.Number,
			Month:  c.Month,
			Year:   c.Year,
			CVC:    c.CVC,
		},
		Desc: "Apartment Reservation",
	}
	ch, err := charge.New(params)
	if err != nil {
		log.Print(err)
		return err
	}
	log.Print(ch)
	return nil
}

type creditCard struct {
	Number string
	Month  string
	Year   string
	CVC    string
}

const minAmountCents = 100    // $1.00
const maxAmountCents = 100000 // $1,000.00

type amount struct {
	Cents int // positive or negative
}

func (a amount) InDollars() string {
	return fmt.Sprintf("$%.02f", float32(a.Cents)/100)
}

type entry struct {
	Amount amount
	Date   time.Time
}

type rateCode int

const (
	rateWithBunny rateCode = iota
	rateWithoutBunny
)

func (c rateCode) Name() string {
	switch c {
	case rateWithoutBunny:
		return "Without Bunny"
	case rateWithBunny:
		return "With Bunnys"
	default:
		panic("unknown rate code")
	}
}

func (c rateCode) String() string {
	return c.Name()
}

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
