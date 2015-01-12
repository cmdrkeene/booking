package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"strings"
	"text/template"
	"time"

	"github.com/cmdrkeene/booking"
)

var flagDataPath = flag.String("data", "/tmp/booking.db", "Path to data file")
var flagHttp = flag.String("http", ":3000", "HTTP address e.g. :3000")
var flagProcessorToken = flag.String("processor-token", "test", "Payment processor token")

func main() {
	flag.Parse()

	service := booking.NewService(
		*flagDataPath,
		*flagProcessorToken,
	)

	fmt.Println("listening on", *flagHttp)
	app := newBookingApp(service)
	log.Fatal(app.ListenAndServe(*flagHttp))
}

/*
A web app for guests to register and pay for reservations.

Pages:
* Home - shows a calendar of availability
* Book - date picker, if unavailable show error and repeat
* Register - after date is picked, register as a guest
* Pay - after dates picked, guest registered, charge credit card
* Complete - after credit card charged, confirm reservation

Store session information in a "session" cookie mapped to an in memory store.

Dependencies:
- way to show availability
- way to create a booking
*/
type bookingApp struct {
	service booking.Service
	mux     http.Handler
}

const (
	pathHome     = "/"
	pathDates    = "/dates"
	pathRegister = "/register"
	pathPay      = "/pay"
	pathComplete = "/complete"
)

const (
	templateHome     = "app/guest/home"
	templateDates    = "app/guest/dates"
	templateRegister = "app/guest/register"
	templatePay      = "app/guest/pay"
	templateComplete = "app/guest/complete"
)

func newBookingApp(service booking.Service) bookingApp {
	app := bookingApp{}
	app.service = service

	mux := http.NewServeMux()
	mux.HandleFunc(pathHome, app.Home)
	mux.HandleFunc(pathDates, app.ChooseDates)
	mux.HandleFunc(pathRegister, app.register)
	mux.HandleFunc(pathPay, app.pay)
	mux.HandleFunc(pathComplete, app.complete)
	app.mux = mux

	return app
}

func (a bookingApp) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	a.mux.ServeHTTP(w, r)
}

func (a bookingApp) ListenAndServe(addr string) error {
	s := http.Server{}
	s.Addr = addr
	s.Handler = a.mux
	return s.ListenAndServe()
}

var prettyDateFormat = "January 2, 2006"
var shortDateFormat = "2006-01-02"

func (a bookingApp) Home(w http.ResponseWriter, r *http.Request) {
	days, err := a.service.AvailableDays()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var vars = struct {
		Days   []time.Time
		Action string
	}{days, pathDates}
	renderTemplate(w, templateHome, vars)
}

func (a bookingApp) ChooseDates(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Not Implemented", http.StatusNotImplemented)
		return
	}

	err := r.ParseForm()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	for _, v := range r.Form["dates"] {
		t, err := time.Parse(shortDateFormat, v)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		w.Write([]byte(t.Format(shortDateFormat) + "\n"))
	}

	return
}
func (a bookingApp) register(w http.ResponseWriter, r *http.Request) {}
func (a bookingApp) pay(w http.ResponseWriter, r *http.Request)      {}
func (a bookingApp) complete(w http.ResponseWriter, r *http.Request) {}

func prettyDate(t time.Time) string {
	return t.Format(prettyDateFormat)
}

func shortDate(t time.Time) string {
	return t.Format(shortDateFormat)
}

func renderTemplate(w http.ResponseWriter, tmpl string, data interface{}) {
	funcs := template.FuncMap{
		"prettyDate": prettyDate,
		"shortDate":  shortDate,
	}

	parts := strings.Split(tmpl, "/")
	name := parts[len(parts)-1] + ".html"
	path := tmpl + ".html"

	t, err := template.New(name).Funcs(funcs).ParseFiles(path)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = t.Execute(w, data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
