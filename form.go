package booking

import (
	"errors"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// Create a new booking from an HTML form
type bookingForm struct {
	Name       string
	Email      string
	Phone      string
	Dates      []time.Time
	CardNumber string
	CardMonth  string
	CardYear   string
	CardCVC    string
}

func newBookingForm(r *http.Request) (bookingForm, error) {
	f := bookingForm{}
	f.Name = r.FormValue("name")
	f.Email = r.FormValue("email")
	f.Phone = r.FormValue("phone")
	f.Phone = r.FormValue("phone")
	for _, v := range r.Form["dates"] {
		t, err := time.Parse(iso8601, v)
		if err != nil {
			return bookingForm{}, datesInvalid
		}
		f.Dates = append(f.Dates, t)
	}
	f.CardNumber = r.FormValue("card_number")
	f.CardMonth = r.FormValue("card_month")
	f.CardYear = r.FormValue("card_year")
	f.CardCVC = r.FormValue("card_cvc")

	return f, nil
}

func (f bookingForm) Encode() string {
	v := &url.Values{}
	v.Set("name", f.Name)
	v.Set("email", f.Email)
	v.Set("phone", f.Phone)
	for _, t := range f.Dates {
		v.Add("dates", t.Format(iso8601))
	}
	v.Set("card_number", f.CardNumber)
	v.Set("card_month", f.CardMonth)
	v.Set("card_year", f.CardYear)
	v.Set("card_cvc", f.CardCVC)
	return v.Encode()
}

func (f bookingForm) Validate() error {
	if f.Name == "" {
		return nameMissing
	}

	if f.Email == "" {
		return emailMissing
	}

	if !strings.Contains(f.Email, "@") {
		return emailInvalid
	}

	if f.Phone == "" {
		return phoneMissing
	}

	if len(f.Dates) == 0 {
		return datesMissing
	}

	if f.CardNumber == "" {
		return cardNumberMissing
	}

	if f.CardMonth == "" {
		return cardMonthMissing
	}

	if f.CardYear == "" {
		return cardYearMissing
	}

	if f.CardCVC == "" {
		return cardCVCMissing
	}

	return nil
}

var phoneMissing = errors.New("phone number missing")
var datesMissing = errors.New("dates missing")
var datesInvalid = errors.New("dates not in YYYY-MM-DD format")
var cardNumberMissing = errors.New("card number missing")
var cardYearMissing = errors.New("card expiration year missing")
var cardMonthMissing = errors.New("card expiration month missing")
var cardCVCMissing = errors.New("card CVC missing")
