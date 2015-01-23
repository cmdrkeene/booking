package booking

import (
	"net/http"
	"net/url"
	"strings"
	"testing"

	"github.com/facebookgo/inject"
)

func TestServer(t *testing.T) {
	// mostly routing, excercise
}

func TestController(t *testing.T) {
	// GET
	// POST
}

func TestForm(t *testing.T) {
	db := testDB()
	defer db.Close()
	var calendar Calendar
	var form Form
	var guestbook Guestbook
	var ledger Ledger
	var register Register
	err := inject.Populate(
		db,
		&calendar,
		&form,
		&guestbook,
		&ledger,
		&register,
	)
	if err != nil {
		t.Error(err)
	}

	// set availability
	calendar.Add(
		newDate(2015, 1, 1),
		newDate(2015, 1, 2),
	)

	// post form
	vals := url.Values{}
	vals.Add(formKeyCardCVC, "123")
	vals.Add(formKeyCardMonth, "01")
	vals.Add(formKeyCardNumber, "1111222233334444")
	vals.Add(formKeyCardYear, "15")
	vals.Add(formKeyCheckin, "2015-01-01")
	vals.Add(formKeyCheckout, "2015-01-02")
	vals.Add(formKeyEmail, "a@b")
	vals.Add(formKeyName, "a b")

	req, err := http.NewRequest("POST", "/", strings.NewReader(vals.Encode()))
	if err != nil {
		t.Error(err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; param=value")

	_, errs := form.Submit(req)
	if len(errs) > 0 {
		t.Error(errs)
	}
}
