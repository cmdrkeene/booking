package booking

import (
	"bytes"
	"database/sql"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"code.google.com/p/go.net/html"
)

func TestBookingControllerNew(t *testing.T) {
	db, _ := sql.Open("sqlite3", ":memory:")
	defer db.Close()

	controller := newBookingController(db)

	// fake availability in test
	av := &testAvailablity{}
	av.Add(time.Date(2015, 2, 1, 0, 0, 0, 0, time.UTC))
	av.Add(time.Date(2015, 2, 2, 0, 0, 0, 0, time.UTC))
	controller.availability = av

	w := httptest.NewRecorder()
	r, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Error(err)
	}
	controller.ServeHTTP(w, r)

	// check 200 OK
	if w.Code != http.StatusOK {
		t.Error("want", http.StatusOK)
		t.Error("got ", w.Code)
	}

	// check availability shown
	mustContain(t, w.Body, "February 1, 2015")
	mustContain(t, w.Body, "February 2, 2015")

	// check rates shown
	mustContain(t, w.Body, "With Bunny")
	mustContain(t, w.Body, "Without Bunny")

	// check form
	// TODO make more rigorous - perhaps bookingForm can render its own template?
	doc, err := html.Parse(w.Body)
	if err != nil {
		t.Error(err)
	}
	forms := getElements(doc, "form")
	if len(forms) == 0 {
		t.Error("want html form")
		t.Error("got none")
	}
}

func TestBookingControllerCreate(t *testing.T) {
	emptyForm := bookingForm{}
	completeForm := bookingForm{
		Name:  "Brandon",
		Email: "a@b.com",
		Phone: "(555) 111-1212",
		Dates: []time.Time{
			time.Date(2015, 2, 1, 0, 0, 0, 0, time.UTC),
			time.Date(2015, 2, 2, 0, 0, 0, 0, time.UTC),
		},
		CardNumber: "1111222233334444",
		CardMonth:  "1",
		CardYear:   "2015",
		CardCVC:    "111",
	}

	var tests = []struct {
		form bookingForm
		code int
		err  error
	}{
		{emptyForm, http.StatusBadRequest, nameMissing},
		{completeForm, http.StatusCreated, nil},
	}

	db, _ := sql.Open("sqlite3", ":memory:")
	defer db.Close()

	controller := newBookingController(db)

	for _, tt := range tests {
		url := "/?" + tt.form.Encode()
		r, err := http.NewRequest("POST", url, nil)
		if err != nil {
			t.Error(err)
		}

		w := httptest.NewRecorder()
		controller.ServeHTTP(w, r)

		if w.Code != tt.code {
			t.Error("want", tt.code)
			t.Error("got ", w.Code)
		}

		if tt.err != nil {
			if !bytes.Contains(w.Body.Bytes(), []byte(tt.err.Error())) {
				t.Error("want", tt.err.Error())
				t.Error("got ", string(w.Body.Bytes()))
			}
		}
	}
}

// helper for checking if strings occur in Buffer
func mustContain(t *testing.T, b *bytes.Buffer, s string) {
	if !bytes.Contains(b.Bytes(), []byte(s)) {
		t.Error("want", s)
		t.Error("got ", string(b.Bytes()))
	}
}

func getElements(doc *html.Node, name string) []*html.Node {
	var found []*html.Node
	var f func(*html.Node)
	f = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == name {
			found = append(found, n)
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}
	f(doc)
	return found
}
