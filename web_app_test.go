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

func TestWebApp(t *testing.T) {
	service := testService{
		availableDays: []time.Time{
			time.Date(2015, 2, 1, 0, 0, 0, 0, time.UTC),
			time.Date(2015, 2, 2, 0, 0, 0, 0, time.UTC),
		},
	}

	app := newWebApp(service)
	ts := httptest.NewServer(app)
	defer ts.Close()

	// new test client
	client := newTestClient(ts.URL)

	// get homepage
	client.Get("/")
	if code := client.response.StatusCode; code != http.StatusOK {
		t.Error("want", http.StatusOK)
		t.Error("got ", code)
	}

	// ensure dates are listed
	want := []byte("February 1, 2015")
	if !bytes.Contains(client.body, want) {
		t.Error("want contains", string(want))
		t.Error("got", string(client.body))
	}

	want = []byte("February 2, 2015")
	if !bytes.Contains(client.body, want) {
		t.Error("want contains", string(want))
		t.Error("got", string(client.body))
	}

	// find form
	doc, err := html.Parse(bytes.NewReader(client.body))
	if err != nil {
		t.Error(err)
	}

	forms := getElements(doc, "form")
	if len(forms) == 0 {
		t.Error("want html form")
		t.Error("got none")
	}

	// select no checkboxes
	// select two checkboxes
	// click Book
	// see registration
}

type testClient struct {
	url      string
	cookies  *http.CookieJar
	visited  []string       // list of visited urls. head is first, tail is last.
	response *http.Response // last response
	body     []byte
	code     int
}

func newTestClient(url string) *testClient {
	c := &testClient{}
	c.url = url
	return c
}

func (c *testClient) Response() *http.Response {
	return c.response
}

func (c *testClient) Get(path string) error {
	var err error
	c.response, err = http.Get(c.url + path)
	if err != nil {
		c.body = nil
		c.code = 0
		return err
	}
	c.body, err = ioutil.ReadAll(c.response.Body)
	defer c.response.Body.Close()
	if err != nil {
		return err
	}
	c.code = c.response.StatusCode

	return nil
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
