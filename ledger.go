package booking

import (
	"database/sql"
	"database/sql/driver"
	"errors"
	"fmt"

	"github.com/golang/glog"
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

func (r *rate) Scan(src interface{}) error {
	rawName, ok := src.([]byte)
	if !ok {
		err := errors.New(
			fmt.Sprintf("can't scan rate from db: %#v", src),
		)
		glog.Error(err)
		return err
	}
	name := string(rawName)
	for _, element := range allRates {
		if name == element.Name {
			r.Amount = element.Amount
			r.Name = element.Name
			return nil
		}
	}

	err := errors.New("can't find rate for " + name)
	glog.Error(err)
	return err
}

func (r rate) Value() (driver.Value, error) {
	return driver.Value(r.Name), nil
}

// What we charge for the place
var (
	withBunny    = rate{Name: "With Bunny", Amount: amount(20000)}
	withoutBunny = rate{Name: "Without Bunny", Amount: amount(25000)}
	allRates     = []rate{withBunny, withoutBunny}
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
	return 0, nil
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
