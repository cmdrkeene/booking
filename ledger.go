package booking

import (
	"database/sql"
	"fmt"
)

// A sum of money stored in 1/100th increments (e.g. cents)
type amount int64

// Round to nearest cent for simple display
func (a amount) String() string {
	return fmt.Sprintf("$%0.2f", float64(a)/100.0)
}

// A note on a transaction
type memo string

// A named amount to charge for a booking
type rate struct {
	Name   string
	Amount amount
}

// What we charge for the place
var (
	withBunny    = rate{Name: "With Bunny", Amount: amount(20000)}
	withoutBunny = rate{Name: "With Bunny", Amount: amount(25000)}
)

// See: http://www.regular-expressions.info/creditcard.html
type creditCard struct {
	CVC    string
	Month  int
	Number int
	Year   int
}

// Summary of guest accounts
type Ledger struct {
	DB        *sql.DB    `inject:""`
	Guestbook *Guestbook `inject:""`
}

func (l *Ledger) Balance(guestId) (amount, error) {
	return amount(0), nil
}

func (l *Ledger) Debit(guestId, amount, memo) error {
	return nil
}

func (l *Ledger) Credit(guestId, amount, memo) error {
	return nil
}

func (l *Ledger) Charge(guestId, amount, creditCard, memo) error {
	return nil
}
