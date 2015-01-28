package booking

import (
	"database/sql"
	"database/sql/driver"
	"errors"
	"fmt"
	"log"

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

func (r rate) String() string {
	return fmt.Sprintf("rate: %s (%s)", r.Name, r.Amount.String())
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
	withBunny = rate{
		Amount: amount(20000),
		Name:   "With Bunny",
	}
	withoutBunny = rate{
		Amount: amount(25000),
		Name:   "Without Bunny",
	}
	allRates = []rate{withBunny, withoutBunny}
)

// limits "new" rates a predefined list
func newRate(name string) (rate, error) {
	for _, r := range allRates {
		if r.Name == name {
			return r, nil
		}
	}
	return rate{}, errors.New("rate not found")
}

// See: http://www.regular-expressions.info/creditcard.html
type creditCard struct {
	CVC    int
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

func (l *Ledger) Charge(
	guest guestId,
	amount amount,
	card creditCard,
	memo memo,
) error {
	log.Println("charged", guest, amount, "on", card)
	return nil
}
