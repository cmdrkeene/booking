package booking

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"code.google.com/p/go.net/html"
)

func TestWebAppHome(t *testing.T) {
	app := NewWebApp(testService{
		availableDays: []time.Time{
			time.Date(2015, 2, 1, 0, 0, 0, 0, time.UTC),
			time.Date(2015, 2, 2, 0, 0, 0, 0, time.UTC),
		},
	})
	ts := httptest.NewServer(app)
	defer ts.Close()

	// get homepage
	req, err := http.NewRequest("GET", ts.URL+pathHome, nil)
	if err != nil {
		t.Error(err)
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Error(err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Error("want", http.StatusOK)
		t.Error("got ", resp.StatusCode)
	}

	// ensure dates are listed
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Error(err)
	}
	resp.Body.Close()

	want := []byte("February 1, 2015")
	if !bytes.Contains(body, want) {
		t.Error("want contains", string(want))
		t.Error("got", string(body))
	}

	want = []byte("February 2, 2015")
	if !bytes.Contains(body, want) {
		t.Error("want contains", string(want))
		t.Error("got", string(body))
	}

	// TODO ensure form is correct
	doc, err := html.Parse(bytes.NewReader(body))
	if err != nil {
		t.Error(err)
	}
	forms := getElements(doc, "form")
	if len(forms) == 0 {
		t.Error("want html form")
		t.Error("got none")
	}
}

func TestWebAppBook(t *testing.T) {
	app := NewWebApp(testService{
		availableDays: []time.Time{
			time.Date(2015, 2, 1, 0, 0, 0, 0, time.UTC),
			time.Date(2015, 2, 2, 0, 0, 0, 0, time.UTC),
		},
	})
	ts := httptest.NewServer(app)
	defer ts.Close()

	// submit incomplete form
	form := bookingForm{}
	url := ts.URL + pathBook + "?" + form.Encode()
	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		t.Error(err)
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Error(err)
	}

	// check bad request
	if resp.StatusCode != http.StatusBadRequest {
		t.Error("want", http.StatusTemporaryRedirect)
		t.Error("got ", resp.StatusCode)
	}

	// check error message
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Error(err)
	}
	resp.Body.Close()

	want := []byte("name missing")
	if !bytes.Contains(body, want) {
		t.Error("want contains", string(want))
		t.Error("got", string(body))
	}

	// submit complete form
	form = bookingForm{}
	form.Name = "Brandon"
	form.Email = "a@b.com"
	form.Phone = "(555) 111-1212"
	form.Dates = []time.Time{
		time.Date(2015, 2, 1, 0, 0, 0, 0, time.UTC),
		time.Date(2015, 2, 2, 0, 0, 0, 0, time.UTC),
	}
	form.CardNumber = "1111222233334444"
	form.CardMonth = "1"
	form.CardYear = "2015"
	form.CardCVC = "111"

	url = ts.URL + pathBook + "?" + form.Encode()
	req, err = http.NewRequest("POST", url, nil)
	if err != nil {
		t.Error(err)
	}
	resp, err = http.DefaultClient.Do(req)
	if err != nil {
		t.Error(err)
	}
	resp.Body.Close()

	// check created
	if resp.StatusCode != http.StatusCreated {
		t.Error("want", http.StatusCreated)
		t.Error("got ", resp.StatusCode)
	}
}

type testService struct {
	availableDays []time.Time
}

func (ts testService) AvailableDays() ([]time.Time, error) {
	return ts.availableDays, nil
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
