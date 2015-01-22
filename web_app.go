package booking

import (
	"net/http"

	"github.com/golang/glog"
)

// Register, Charge, and Book in one step from raw input
type Form struct {
	Guestbook *Guestbook `inject:""`
	Ledger    *Ledger    `inject:""`
	Register  *Register  `inject:""`
}

func (form *Form) Submit(r *http.Request) (bookingId, []error) {
	// fields
	fields, errs := newFormFields(r)
	if len(errs) > 0 {
		return 0, errs
	}

	// register
	guestId, err := form.Guestbook.Register(
		fields.Name,
		fields.Email,
		fields.Phone,
	)

	if err != nil {
		glog.Error(err)
		return 0, []error{err}
	}

	// charge
	memo := memo(fields.Rate.String())
	err = form.Ledger.Charge(
		guestId,
		fields.Rate.Amount,
		fields.Card,
		memo,
	)
	if err != nil {
		glog.Error(err)
		return 0, []error{err}
	}

	// book
	bookingId, err := form.Register.Book(
		fields.Checkin,
		fields.Checkout,
		guestId,
		fields.Rate,
	)
	if err != nil {
		glog.Error(err)
		return 0, []error{err}
	}

	return bookingId, []error{}
}

// helper to bucket form fields
type formFields struct {
	Checkin  date
	Checkout date
	Card     creditCard
	Email    email
	Name     name
	Phone    phoneNumber
	Rate     rate
}

// Convert an http.Request to useful value
func newFormFields(r *http.Request) (formFields, []error) {
	errs := []error{}

	var err error
	var fields formFields

	fields.Card, err = newCreditCard(
		r.FormValue("card_cvc"),
		r.FormValue("card_month"),
		r.FormValue("card_number"),
		r.FormValue("card_year"),
	)
	if err != nil {
		errs = append(errs, err)
	}

	fields.Checkin, err = newDateFromString(r.FormValue("checkin"))
	if err != nil {
		errs = append(errs, err)
	}

	fields.Checkout, err = newDateFromString(r.FormValue("checkout"))
	if err != nil {
		errs = append(errs, err)
	}

	fields.Email, err = newEmail(r.FormValue("email"))
	if err != nil {
		errs = append(errs, err)
	}

	fields.Name, err = newName(r.FormValue("name"))
	if err != nil {
		errs = append(errs, err)
	}

	fields.Phone, err = newPhoneNumber(r.FormValue("phone"))
	if err != nil {
		errs = append(errs, err)
	}

	if len(errs) != 0 {
		return formFields{}, errs
	}

	return fields, []error{}
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
