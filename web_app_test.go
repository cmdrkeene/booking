package booking

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/facebookgo/inject"
)

func TestServer(t *testing.T) {
	// mostly routing, excercise
}

func TestControllerGet(t *testing.T) {
	db := testDB()
	defer db.Close()
	var calendar Calendar
	var controller Controller
	err := inject.Populate(
		db,
		&calendar,
		&controller,
	)
	if err != nil {
		t.Error(err)
	}
	calendar.Add(
		newDate(2015, 1, 1),
		newDate(2015, 1, 2),
		newDate(2015, 1, 3),
	)

	w := httptest.NewRecorder()
	r, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Error(err)
	}
	controller.Get(w, r)
	if w.Code != http.StatusOK {
		t.Error("want", http.StatusOK)
		t.Error("got ", w.Code)
	}

	body := w.Body.Bytes()
	dates := []string{
		"2015-01-01",
		"2015-01-02",
		"2015-01-03",
	}
	for _, s := range dates {
		if !bytes.Contains(body, []byte(s)) {
			t.Error("want", s)
			t.Error("got ", string(body))
		}
	}
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
	req.Header.Set(
		"Content-Type",
		"application/x-www-form-urlencoded; param=value",
	)

	_, errs := form.Submit(req)
	if len(errs) > 0 {
		t.Error(errs)
	}
}
