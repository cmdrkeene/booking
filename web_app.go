package booking

import (
	"html/template"
	"net/http"
	"time"
)

type webApp struct {
	service   Service
	mux       *http.ServeMux
	templates map[string]*template.Template
}

const templateHome = `
<html>
  <body>
    <h1>15 Dunham Place</h1>
    <h3>Book your stay</h3>
    <form action="{{.PathBook}}" method="post">
      <fieldset>
        <legend>Dates</legend>
        {{range .Days}}
          <input type="checkbox" name="dates" value="{{ isoDate . }}" />
          {{prettyDate .}}
          <br />
        {{end}}
        </ul>
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
</html>
`

const (
	pathHome = "/"
	pathBook = "/book"
)

func NewWebApp(s Service) webApp {
	app := webApp{}
	app.service = s

	app.templates = make(map[string]*template.Template)

	mux := http.NewServeMux()
	mux.HandleFunc(pathHome, app.Home)
	mux.HandleFunc(pathBook, app.Book)
	app.mux = mux

	return app
}

func (a webApp) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	a.mux.ServeHTTP(w, r)
}

func (a webApp) Home(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Not Implemented", http.StatusNotImplemented)
		return
	}

	fm := template.FuncMap{
		"isoDate":    func(t time.Time) string { return t.Format(iso8601) },
		"prettyDate": func(t time.Time) string { return t.Format(pretty) },
	}
	t, err := template.New("home").Funcs(fm).Parse(templateHome)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	days, err := a.service.AvailableDays()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var data = struct {
		Days     []time.Time
		PathBook string
	}{days, pathBook}
	t.Execute(w, data)
}

func (a webApp) Book(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Not Implemented", http.StatusNotImplemented)
		return
	}

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

	w.WriteHeader(http.StatusCreated)
}
