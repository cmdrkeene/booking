package booking

import (
	"net/http"

	"github.com/golang/glog"
)

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
