package booking

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/cmdrkeene/booking/pkg/date"
)

// Builds Form instances
type FormBuilder struct {
	Calendar  *Calendar  `inject:""`
	Guestbook *Guestbook `inject:""`
	Ledger    *Ledger    `inject:""`
	Register  *Register  `inject:""`
}

func (b *FormBuilder) Build() *Form {
	return &Form{
		calendar:  b.Calendar,
		guestbook: b.Guestbook,
		ledger:    b.Ledger,
		register:  b.Register,
	}
}

// Register, Charge, and Book in one-step
// Validates input state
type Form struct {
	// Dependencies
	calendar  *Calendar
	guestbook *Guestbook
	ledger    *Ledger
	register  *Register

	// Errors
	Errors map[string]string

	// Fields
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

// Converts http.Request to Form and batches errors
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

	fmt.Printf("%#v", form)

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
	validator.Integer(fvCardCVC, form.CardCVC)
	validator.Integer(fvCardMonth, form.CardMonth)
	validator.Integer(fvCardNumber, form.CardNumber)
	validator.Integer(fvCardYear, form.CardYear)
	validator.Date(fvCheckin, form.Checkin)
	validator.Date(fvCheckout, form.Checkout)
	validator.Email(fvEmail, form.Email)
	validator.Name(fvName, form.Name)
	validator.Phone(fvPhone, form.Phone)
	validator.Rate(fvRate, form.Rate)

	// do not continue if anything is invalid
	if len(validator.Errors) > 0 {
		form.Errors = validator.Errors
		return false
	}

	return true
}

// func (form *Form) xSubmit(r *http.Request) (bookingId, []error) {
// 	// register
// 	guestId, err := form.Guestbook.Register(
// 		fields.Name,
// 		fields.Email,
// 		fields.Phone,
// 	)
// 	if err != nil {
// 		glog.Error(err)
// 		return 0, []error{err}
// 	}

// 	// charge
// 	memo := memo(fields.Rate.String())
// 	err = form.Ledger.Charge(
// 		guestId,
// 		fields.Rate.Amount,
// 		fields.Card,
// 		memo,
// 	)
// 	if err != nil {
// 		glog.Error(err)
// 		return 0, []error{err}
// 	}

// 	// book
// 	bookingId, err := form.Register.Book(
// 		fields.Checkin,
// 		fields.Checkout,
// 		guestId,
// 		fields.Rate,
// 	)
// 	if err != nil {
// 		glog.Error(err)
// 		return 0, []error{err}
// 	}

// 	return bookingId, []error{}
// }

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

func (val validator) Integer(k, v string) {
	if _, err := strconv.Atoi(v); err != nil {
		val.Errors[k] = "invalid"
	}
}

func (val validator) Date(k, v string) {
	log.Print("validating", v)
	if _, err := date.Parse(v); err != nil {
		val.Errors[k] = "invalid"
	}
}

func (val validator) Name(k, v string) {
	if _, err := newName(v); err != nil {
		val.Errors[k] = "invalid"
	}
}

func (val validator) Email(k, v string) {
	if _, err := newEmail(v); err != nil {
		val.Errors[k] = "invalid"
	}
}

func (val validator) Phone(k, v string) {
	if _, err := newPhoneNumber(v); err != nil {
		val.Errors[k] = "invalid"
	}
}

func (val validator) Rate(k, v string) {
	if _, err := newRate(v); err != nil {
		val.Errors[k] = "invalid"
	}
}
