package booking

import (
	"bytes"
	"fmt"
	"html/template"
	"log"

	"net/http"
	"strings"

	"github.com/golang/glog"
)

// Handle HTTP interaction with Form
type Handler struct {
	Calendar *Calendar `inject:""`
	Form     *Form     `inject:""`
}

var templateHelpers = template.FuncMap{
	"formatDate": func(d date) string { return d.Format(layoutDatePretty) },
}

// Display form with available dates listed
func (h *Handler) Get(w http.ResponseWriter, r *http.Request) {
	t := template.New("form.html.go")
	t = t.Funcs(templateHelpers)
	t = template.Must(t.Parse(formHtml))

	available, err := h.Calendar.List()
	if err != nil {
		glog.Error(err)
		http.Error(w, err.Error(), 500)
		return
	}

	var data = struct {
		AvailableDates []date
		Rates          []rate
	}{available, allRates}
	h.writeTemplate(w, t, data)
}

// Submit form and display errors or confirmation page
func (h *Handler) Post(w http.ResponseWriter, r *http.Request) {
	bookingId, errors := h.Form.Submit(r)
	if len(errors) > 0 {
		var errorMessages []string
		for _, err := range errors {
			errorMessages = append(errorMessages, err.Error())
		}
		s := strings.Join(errorMessages, ", ")
		http.Error(w, s, 500)
		return
	}

	w.WriteHeader(201)
	message := fmt.Sprintf("success! booking confirmation id %d", bookingId)
	w.Write([]byte(message))
}

/*
Booking HTTP Handler

GET /   -> 200 lists availability and form
POST /  -> 201 submits form (register, pay, and book workflow)

*/
func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		h.Get(w, r)
		return
	}

	if r.Method == "POST" {
		log.Print("calling post")
		h.Post(w, r)
		return
	}

	http.Error(w, "Not Implemented", http.StatusNotImplemented)
}

// buffers template Execute and returns 500 on any error
// avoids returning 200 OK when Execute fails after starting Write()
func (h *Handler) writeTemplate(
	w http.ResponseWriter,
	t *template.Template,
	data interface{},
) {
	var b bytes.Buffer
	err := t.Execute(&b, data)
	if err != nil {
		glog.Error(err)
		http.Error(w, http.StatusText(500), 500)
	}
	b.WriteTo(w)
}

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

const (
	formKeyCardCVC    = "card_cvc"
	formKeyCardMonth  = "card_month"
	formKeyCardNumber = "card_number"
	formKeyCardYear   = "card_year"
	formKeyCheckin    = "checkin"
	formKeyCheckout   = "checkout"
	formKeyEmail      = "email"
	formKeyName       = "name"
	formKeyPhone      = "phone"
)

// Convert an http.Request to useful value
func newFormFields(r *http.Request) (formFields, []error) {
	errs := []error{}

	var err error
	var fields formFields

	fields.Card, err = newCreditCard(
		r.FormValue(formKeyCardCVC),
		r.FormValue(formKeyCardMonth),
		r.FormValue(formKeyCardNumber),
		r.FormValue(formKeyCardYear),
	)
	if err != nil {
		errs = append(errs, err)
	}

	fields.Checkin, err = newDateFromString(r.FormValue(formKeyCheckin))
	if err != nil {
		errs = append(errs, err)
	}

	fields.Checkout, err = newDateFromString(r.FormValue(formKeyCheckout))
	if err != nil {
		errs = append(errs, err)
	}

	fields.Email, err = newEmail(r.FormValue(formKeyEmail))
	if err != nil {
		errs = append(errs, err)
	}

	fields.Name, err = newName(r.FormValue(formKeyName))
	if err != nil {
		errs = append(errs, err)
	}

	fields.Phone, err = newPhoneNumber(r.FormValue(formKeyPhone))
	if err != nil {
		errs = append(errs, err)
	}

	if len(errs) != 0 {
		return formFields{}, errs
	}

	return fields, []error{}
}
