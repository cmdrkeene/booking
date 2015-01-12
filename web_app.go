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
<ul>
{{range .Days}}
  <li>{{shortDate .}}</li>
{{end}}
</ul>
</body>
</html>
`

const (
	pathHome = "/"
)

func newWebApp(s Service) webApp {
	app := webApp{}
	app.service = s

	app.templates = make(map[string]*template.Template)

	mux := http.NewServeMux()
	mux.HandleFunc(pathHome, app.Home)
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

	fm := template.FuncMap{"shortDate": shortDate}
	t, err := template.New("home").Funcs(fm).Parse(templateHome)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var data = struct {
		Days []time.Time
	}{
		[]time.Time{
			time.Date(2015, 2, 1, 0, 0, 0, 0, time.UTC),
			time.Date(2015, 2, 2, 0, 0, 0, 0, time.UTC),
		},
	}
	t.Execute(w, data)
}

const shortDateFormat = "January 2, 2006"

func shortDate(t time.Time) string {
	return t.Format(shortDateFormat)
}
