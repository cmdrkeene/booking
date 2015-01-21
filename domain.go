/*
Book and pay for our apartment during specific available dates.

To make a booking

* Dates must be available

* Dates must be contiguous

* Guest must be registered

* CheckIn must be before CheckOut

* CheckOut must be at least one day after CheckIn

* Rate must be debited from ledger

To make a payment

* Amount must be non zero

* CreditCard must be charged for amount

* Credit must be recorded in ledger

*/
package booking

import (
	"database/sql"
	"net/http"
	"time"

	"github.com/golang/glog"
)

// List of dates available for booking
type Calendar struct {
	DB *sql.DB `inject:""`
}

// Simple date
type date time.Time

// Init creates table if needed
func (c *Calendar) Init() {
	_, err := c.DB.Exec(`
		create table calendar (
			date datetime,
			available bool default true
		)
	`)

	if err == nil {
		glog.Info("calendar table created")
	} else {
		glog.Warning(err)
	}
}

func (c *Calendar) Add(date) error {
	return nil
}

func (c *Calendar) Remove(date) error {
	return nil
}

func (c *Calendar) List() ([]date, error) {
	return []date{}, nil
}

// Locator for a booking record
type bookingId uint32

// Official list of bookings
type Register struct {
	Calendar *Calendar `inject:""`
	DB       *sql.DB   `inject:""`
}

func (r *Register) Book(
	checkIn date,
	checkOut date,
	guest guestId,
	rate rate,
) (bookingId, error) {
	return bookingId(0), nil
}

func (r *Register) Cancel(bookingId) error {
	return nil
}

// A sum of money stored in 1/100th increments (e.g. cents)
type amount uint32

// A note on a transaction
type memo string

// A named amount to charge for a booking
type rate int

// See: http://www.regular-expressions.info/creditcard.html
type creditCard struct {
	CVC    string
	Month  int
	Number int
	Year   int
}

// Summary of guest accounts
type Ledger struct {
	DB        *sql.DB    `inject:""`
	Guestbook *Guestbook `inject:""`
}

func (l *Ledger) Balance(guestId) (amount, error) {
	return amount(0), nil
}

func (l *Ledger) Debit(guestId, amount, memo) error {
	return nil
}

func (l *Ledger) Credit(guestId, amount, memo) error {
	return nil
}

func (l *Ledger) Charge(guestId, amount, creditCard, memo) error {
	return nil
}

// Locator for a guest record
type guestId uint32

// A guest's name - must be 2 parts, not empty
type guestName string

// An electronic mail address
type email string

// A telephone number
type phoneNumber string

// A guest record
type guest struct {
	Email       email
	Id          guestId
	Name        guestName
	PhoneNumber phoneNumber
}

// List of registered guests
type Guestbook struct {
	DB *sql.DB `inject:""`
}

func (g *Guestbook) IsRegistered(guestId) (bool, error) {
	return false, nil
}

func (g *Guestbook) Lookup(guestId) (guest, error) {
	return guest{}, nil
}

func (g *Guestbook) Register(guestName, email, phoneNumber) (guestId, error) {
	return guestId(0), nil
}

// Register, Book, and Charge in one step from raw input
type Form struct {
	Guestbook *Guestbook `inject:""`
	Ledger    *Ledger    `inject:""`
	Register  *Register  `inject:""`
}

func (f *Form) Submit(FormFields) (bookingId, error) {
	// TODO maybe use bookingFormSubmission to validate and aggregate errors
	return bookingId(0), nil
}

// A bucket for raw user input
type FormFields struct {
	Name             string
	Email            string
	Phone            string
	CheckIn          string
	CheckOut         string
	CreditCardNumber string
	CreditCardMonth  string
	CreditCardYear   string
	CreditCardCVC    string
}

// Check input and return any errors
func (f *FormFields) Validate() []error {
	return []error{}
}

// Handle HTTP interaction with Form
type Controller struct {
	Calendar *Calendar `inject:""`
	Form     *Form     `inject:""`
}

// Display form with available dates listed
func (c *Controller) Get(w http.ResponseWriter, r *http.Request) {
	glog.Infoln("GET", r.RequestURI)
}

// Submit form and display errors or confirmation page
func (c *Controller) Post(w http.ResponseWriter, r *http.Request) {
	glog.Infoln("POST", r.RequestURI)
}

// Serve Controller
type Server struct {
	Controller *Controller `inject:""`
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		s.Controller.Get(w, r)
		return
	}

	if r.Method == "POST" {
		s.Controller.Post(w, r)
		return
	}

	http.Error(w, "Not Implemented", http.StatusNotImplemented)
}

func (s *Server) ListenAndServe(addr string) error {
	glog.Infoln("listening on", addr)
	return http.ListenAndServe(addr, s)
}
