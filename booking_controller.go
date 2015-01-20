package booking

import (
	"html/template"
	"net/http"
	"time"
)

// Web UI for making a one-step booking
type bookingController struct {
	availability timeLister
	reservations reserver
	template     *template.Template
	path         string
}

type timeLister interface {
	List() ([]time.Time, error)
}

func newBookingController(reg *Registry) bookingController {
	controller := bookingController{}
	controller.availability = newAvailabilityTable(reg.DB())
	controller.path = "/"
	controller.template = newBookingTemplate()
	return controller
}

var dateHelpers = template.FuncMap{
	"isoDate":    func(t time.Time) string { return t.Format(iso8601) },
	"prettyDate": func(t time.Time) string { return t.Format(pretty) },
}

func (c bookingController) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// RESTful routes for path
	if r.Method == "GET" {
		c.new(w, r)
		return
	}

	if r.Method == "POST" {
		c.create(w, r)
		return
	}

	http.Error(w, "Not Implemented", http.StatusNotImplemented)
}

func (c bookingController) new(w http.ResponseWriter, r *http.Request) {
	available, err := c.availability.List()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var data = struct {
		Available []time.Time
		Rates     []rateCode
	}{
		available,
		[]rateCode{rateWithBunny, rateWithoutBunny},
	}
	c.template.Execute(w, data)
}

func (c bookingController) create(w http.ResponseWriter, r *http.Request) {
	form, err := newBookingForm(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = form.Validate()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// TODO - figure out who organizes all this jazz
	// err = form.Submit()
	// if err != nil {
	// 	http.Error(w, err.Error(), http.StatusInternalServerError)
	// 	return
	// }

	w.WriteHeader(http.StatusCreated)
}

func newBookingTemplate() *template.Template {
	t := template.New("booking.new")
	t = t.Funcs(dateHelpers)
	src := `
<html>
  <body>
    <h1>15 Dunham Place</h1>
    <h3>Book your stay</h3>
    <form action="/" method="post">
      <fieldset>
        <legend>Dates</legend>
        {{range .Available}}
          <input type="checkbox" name="dates" value="{{ isoDate . }}" />
          {{prettyDate .}}
          <br />
        {{end}}
        </ul>
      </fieldset>
      <fieldset>
        <legend>Rate</legend>
        {{range .Rates}}
        <div>
          <input name="rate" type="radio" value="WithBunny" />
          <b>{{.Amount.InDollars}}</b>
          {{.Name}}
        </div>
        {{end}}
      </fieldset>
      <fieldset>
        <legend>Contact</legend>
        <table>
          <tr>
            <th>Name</th>
            <td><input type="text" name="name" /></td>
          </tr>
          <tr>
            <th>Email</th>
            <td><input type="text" name="email" /></td>
          </tr>
          <tr>
            <th>Phone</th>
            <td><input type="text" name="phone" /></td>
          </tr>
        </table>
      </fieldset>
      <fieldset>
        <legend>Credit Card</legend>
        <table>
          <tr>
            <th colspan="3" align="left">Number</th>
          </tr>
          <tr>
            <td colspan="3">
              <input type="text" name="card_number" />
            </td>
          </tr>
          <tr>
            <th>Month</th>
            <th>Year</th>
            <th>CVC</th>
          </tr>
          <tr>
            <td>
              <input type="text" name="card_month" size="4" />
            </td>
            <td>
              <input type="text" name="card_year" size="4"/>
            </td>
            <td>
              <input type="password" name="card_cvc" size="4"/>
            </td>
          </tr>
        </table>
      </fieldset>
      <input type="submit" value="Book" />
    </form>
  </body>
</html>`
	t, err := t.Parse(src)
	if err != nil {
		panic(err)
	}
	return t
}
