package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"text/template"

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
	mux.HandleFunc(pathHome, app.home)
	mux.HandleFunc(pathDates, app.dates)
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

func (a bookingApp) home(w http.ResponseWriter, r *http.Request) {
	days, err := a.service.AvailableDays()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	dayStrings := make([]string, len(days))
	for i, d := range days {
		dayStrings[i] = d.Format("January 2, 2006")
	}

	var vars = struct {
		Days   []string
		Action string
	}{dayStrings, pathDates}
	renderTemplate(w, templateHome, vars)
}

func (a bookingApp) dates(w http.ResponseWriter, r *http.Request)    {}
func (a bookingApp) register(w http.ResponseWriter, r *http.Request) {}
func (a bookingApp) pay(w http.ResponseWriter, r *http.Request)      {}
func (a bookingApp) complete(w http.ResponseWriter, r *http.Request) {}

func renderTemplate(w http.ResponseWriter, tmpl string, data interface{}) {
	t, err := template.ParseFiles(tmpl + ".html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = t.Execute(w, data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
