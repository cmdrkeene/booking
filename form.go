package booking

import (
	"database/sql"
	"net/http"
	"strconv"

	"github.com/cmdrkeene/booking/pkg/date"
	"github.com/golang/glog"
)

// Builds Form instances
type FormBuilder struct {
	Calendar *Calendar `inject:""`
	DB       *sql.DB   `inject:""`
	Ledger   *Ledger   `inject:""`
	Register *Register `inject:""`
}

func (b *FormBuilder) Build() *Form {
	return &Form{
		calendar: b.Calendar,
		db:       b.DB,
		ledger:   b.Ledger,
		register: b.Register,
	}
}

// One page booking form
type Form struct {
	// Dependencies
	calendar *Calendar
	db       *sql.DB
	ledger   *Ledger
	register *Register

	// Errors
	Errors map[string]string

	// Input
	CardCVC    string
	CardMonth  string
	CardNumber string
	CardYear   string
	Checkin    string
	Checkout   string
	Email      string
	Name       string
	Phone      string
	Rate       string

	// Private, valid fields
	cardCVC    int
	cardMonth  int
	cardNumber int
	cardYear   int
	checkin    date.Date
	checkout   date.Date
	email      email
	name       name
	phone      phoneNumber
	rate       rate
}

func (form *Form) AvailableDates() []date.Date {
	dates, err := form.calendar.List()
	if err != nil {
		panic(err)
	}
	return dates
}

func (form *Form) Rates() []rate {
	return allRates
}

const (
	fvCardCVC    = "CardCVC"
	fvCardMonth  = "CardMonth"
	fvCardNumber = "CardNumber"
	fvCardYear   = "CardYear"
	fvCheckin    = "Checkin"
	fvCheckout   = "Checkout"
	fvEmail      = "Email"
	fvName       = "Name"
	fvPhone      = "Phone"
	fvRate       = "Rate"
)

// Converts http.Request to Form. Caller must check Errors
func (form *Form) Validate(r *http.Request) bool {
	form.Errors = make(map[string]string)

	// map request to fields
	form.CardCVC = r.FormValue(fvCardCVC)
	form.CardMonth = r.FormValue(fvCardMonth)
	form.CardNumber = r.FormValue(fvCardNumber)
	form.CardYear = r.FormValue(fvCardYear)
	form.Checkin = r.FormValue(fvCheckin)
	form.Checkout = r.FormValue(fvCheckout)
	form.Email = r.FormValue(fvEmail)
	form.Name = r.FormValue(fvName)
	form.Phone = r.FormValue(fvPhone)
	form.Rate = r.FormValue(fvRate)

	// new validator
	validator := newValidator()

	// check required
	validator.Require(fvCardCVC, form.CardCVC)
	validator.Require(fvCardMonth, form.CardMonth)
	validator.Require(fvCardNumber, form.CardNumber)
	validator.Require(fvCardYear, form.CardYear)
	validator.Require(fvCheckin, form.Checkin)
	validator.Require(fvCheckout, form.Checkout)
	validator.Require(fvEmail, form.Email)
	validator.Require(fvName, form.Name)
	validator.Require(fvPhone, form.Phone)
	validator.Require(fvRate, form.Rate)

	// do not continue if anything is missing
	if len(validator.Errors) > 0 {
		form.Errors = validator.Errors
		return false
	}

	// check data format
	form.cardCVC = validator.Integer(fvCardCVC, form.CardCVC)
	form.cardMonth = validator.Integer(fvCardMonth, form.CardMonth)
	form.cardNumber = validator.Integer(fvCardNumber, form.CardNumber)
	form.cardYear = validator.Integer(fvCardYear, form.CardYear)
	form.checkin = validator.Date(fvCheckin, form.Checkin)
	form.checkout = validator.Date(fvCheckout, form.Checkout)
	form.email = validator.Email(fvEmail, form.Email)
	form.name = validator.Name(fvName, form.Name)
	form.phone = validator.Phone(fvPhone, form.Phone)
	form.rate = validator.Rate(fvRate, form.Rate)

	// do not continue if anything is invalid
	if len(validator.Errors) > 0 {
		form.Errors = validator.Errors
		return false
	}

	return true
}

// Validate, Register, Charge, and Book in one-step
func (form *Form) Submit(r *http.Request) (bookingId, bool) {
	// validate
	if !form.Validate(r) {
		return 0, false
	}

	// start a transaction
	tx, err := form.db.Begin()
	if err != nil {
		panic(err)
	}
	defer tx.Rollback()

	// register guest
	guestbookTx := &GuestbookTx{tx}
	guestId, err := guestbookTx.Register(
		form.name,
		form.email,
		form.phone,
	)
	if err != nil {
		glog.Error(err)
		form.Errors["Register"] = err.Error()
		return 0, false
	}

	// charge card
	err = form.ledger.ChargeTx(
		tx,
		guestId,
		form.rate.Amount,
		form.creditCard(),
		memo(form.rate.String()),
	)
	if err != nil {
		glog.Error(err)
		form.Errors["Charge"] = err.Error()

		return 0, false
	}

	// book dates
	bookingId, err := form.register.BookTx(
		tx,
		form.checkin,
		form.checkout,
		guestId,
		form.rate,
	)
	if err != nil {
		form.Errors["Book"] = err.Error()
		glog.Error(err)
		return 0, false
	}

	// commit and tell user if it failed
	err = tx.Commit()
	if err != nil {
		glog.Error(err)
		form.Errors["Transaction"] = err.Error()
		return 0, false
	}

	return bookingId, true
}

func (f *Form) creditCard() creditCard {
	return creditCard{
		CVC:    f.cardCVC,
		Month:  f.cardMonth,
		Number: f.cardNumber,
		Year:   f.cardYear,
	}
}

type validator struct {
	Errors map[string]string
}

func newValidator() validator {
	return validator{Errors: make(map[string]string)}
}

func (val validator) Require(k, v string) {
	if v == "" {
		val.Errors[k] = "required"
	}
}

func (val validator) Integer(k, v string) int {
	i, err := strconv.Atoi(v)
	if err != nil {
		val.Errors[k] = "invalid"
		return 0
	}
	return i
}

func (val validator) Date(k, v string) date.Date {
	d, err := date.Parse(v)
	if err != nil {
		val.Errors[k] = "invalid"
		return date.Date{}
	}
	return d
}

func (val validator) Name(k, v string) name {
	n, err := newName(v)
	if err != nil {
		val.Errors[k] = "invalid"
		return name{}
	}
	return n
}

func (val validator) Email(k, v string) email {
	e, err := newEmail(v)
	if err != nil {
		val.Errors[k] = "invalid"
		return email{}
	}
	return e
}

func (val validator) Phone(k, v string) phoneNumber {
	p, err := newPhoneNumber(v)
	if err != nil {
		val.Errors[k] = "invalid"
		return phoneNumber{}
	}
	return p
}

func (val validator) Rate(k, v string) rate {
	r, err := newRate(v)
	if err != nil {
		val.Errors[k] = "invalid"
		return rate{}
	}
	return r
}
