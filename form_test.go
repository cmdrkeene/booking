package booking

import (
	"net/http"
	"net/url"
	"reflect"
	"strings"
	"testing"

	"github.com/cmdrkeene/booking/pkg/date"
	"github.com/facebookgo/inject"
)

func TestForm(t *testing.T) {
	db := testDB()
	defer db.Close()
	var cal Calendar
	var formBuilder FormBuilder
	err := inject.Populate(db, &cal, &formBuilder)
	if err != nil {
		t.Error(err)
	}

	// new Form
	form := formBuilder.Build()

	// AvailableDates()
	available := []date.Date{
		date.New(2015, 1, 1),
		date.New(2015, 1, 2),
	}
	cal.Add(available...)
	if got := form.AvailableDates(); !reflect.DeepEqual(available, got) {
		t.Error("want", available)
		t.Error("got ", got)
	}

	// Validate() -> required
	emptyRequest, _ := http.NewRequest("POST", "/", nil)
	ok := form.Validate(emptyRequest)
	if ok {
		t.Error("want Validate() to fail")
	}
	errors := map[string]string{
		"CardCVC":    "required",
		"CardMonth":  "required",
		"CardNumber": "required",
		"CardYear":   "required",
		"Checkin":    "required",
		"Checkout":   "required",
		"Email":      "required",
		"Name":       "required",
		"Phone":      "required",
		"Rate":       "required",
	}
	if !reflect.DeepEqual(errors, form.Errors) {
		t.Error("want", errors)
		t.Error("got ", form.Errors)
	}

	// Validate() -> invalid
	vals := url.Values{}
	vals.Set(fvCardCVC, "a")
	vals.Set(fvCardMonth, "b")
	vals.Set(fvCardNumber, "c")
	vals.Set(fvCardYear, "d")
	vals.Set(fvCheckin, "e")
	vals.Set(fvCheckout, "f")
	vals.Set(fvEmail, "g")
	vals.Set(fvName, "h")
	vals.Set(fvPhone, "i")
	vals.Set(fvRate, "j")

	invalidRequest, _ := http.NewRequest("POST", "/", strings.NewReader(vals.Encode()))
	invalidRequest.Header.Set(
		"Content-Type",
		"application/x-www-form-urlencoded; param=value",
	)
	ok = form.Validate(invalidRequest)
	if ok {
		t.Error("want Validate() to fail")
	}
	errors = map[string]string{
		"CardCVC":    "invalid",
		"CardMonth":  "invalid",
		"CardNumber": "invalid",
		"CardYear":   "invalid",
		"Checkin":    "invalid",
		"Checkout":   "invalid",
		"Email":      "invalid",
		"Name":       "invalid",
		"Phone":      "invalid",
		"Rate":       "invalid",
	}
	if !reflect.DeepEqual(errors, form.Errors) {
		t.Error("want", errors)
		t.Error("got ", form.Errors)
	}

	// Validate() -> ok
	vals = url.Values{}
	vals.Set(fvCardCVC, "123")
	vals.Set(fvCardMonth, "12")
	vals.Set(fvCardNumber, "1111222233334444")
	vals.Set(fvCardYear, "2015")
	vals.Set(fvCheckin, "1/2/2015")
	vals.Set(fvCheckout, "1/3/2015")
	vals.Set(fvEmail, "a@b")
	vals.Set(fvName, "a b")
	vals.Set(fvPhone, "555-123-4567")
	vals.Set(fvRate, withBunny.Name)

	validRequest, _ := http.NewRequest("POST", "/", strings.NewReader(vals.Encode()))
	validRequest.Header.Set(
		"Content-Type",
		"application/x-www-form-urlencoded; param=value",
	)
	ok = form.Validate(validRequest)
	if !ok {
		t.Error("want Validate() to succeed")
	}
	if l := len(form.Errors); l != 0 {
		t.Error("want 0 errors")
		t.Error("got ", form.Errors)
	}
}

// func xTestForm(t *testing.T) {
// 	db := testDB()
// 	defer db.Close()
// 	var calendar Calendar
// 	var form Form
// 	var guestbook Guestbook
// 	var ledger Ledger
// 	var register Register
// 	err := inject.Populate(
// 		db,
// 		&calendar,
// 		&form,
// 		&guestbook,
// 		&ledger,
// 		&register,
// 	)
// 	if err != nil {
// 		t.Error(err)
// 	}

// 	// set availability
// 	calendar.Add(
// 		date.New(2015, 1, 1),
// 		date.New(2015, 1, 2),
// 	)

// 	// post form
// 	vals := url.Values{}
// 	vals.Add(formKeyCardCVC, "123")
// 	vals.Add(formKeyCardMonth, "01")
// 	vals.Add(formKeyCardNumber, "1111222233334444")
// 	vals.Add(formKeyCardYear, "15")
// 	vals.Add(formKeyCheckin, "2015-01-01")
// 	vals.Add(formKeyCheckout, "2015-01-02")
// 	vals.Add(formKeyEmail, "a@b")
// 	vals.Add(formKeyName, "a b")

// 	req, err := http.NewRequest("POST", "/", strings.NewReader(vals.Encode()))
// 	if err != nil {
// 		t.Error(err)
// 	}
// 	req.Header.Set(
// 		"Content-Type",
// 		"application/x-www-form-urlencoded; param=value",
// 	)

// 	_, errs := form.Submit(req)
// 	if len(errs) > 0 {
// 		t.Error(errs)
// 	}
// }
