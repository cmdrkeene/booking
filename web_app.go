package booking

import (
	"bytes"
	"html/template"

	"net/http"

	"github.com/golang/glog"
)

// Handle HTTP interaction with Form
type Handler struct {
	FormBuilder *FormBuilder `inject:""`
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	mux := http.NewServeMux()
	mux.HandleFunc("/", h.index)
	mux.HandleFunc("/confirmation", h.confirmation)
	mux.ServeHTTP(w, r)
}

func (h *Handler) confirmation(w http.ResponseWriter, r *http.Request) {
}

func (h *Handler) index(w http.ResponseWriter, r *http.Request) {
	form := h.FormBuilder.Build()

	switch r.Method {
	case "GET":
		render(w, templateForm, form)
	case "POST":
		if !form.Validate(r) {
			render(w, templateForm, form)
			return
		}

		// bookingId, err := form.Submit()
		// if err != nil {
		// 	http.Error(w, err.Error(), http.StatusInternalServerError)
		// 	return
		// }

		// http.Redirect(w, r, path, http.StatusSeeOther)
	default:
		http.Error(
			w,
			http.StatusText(http.StatusNotImplemented),
			http.StatusNotImplemented,
		)
	}
}

// buffer execute because it can fail midway
func render(w http.ResponseWriter, tmpl *template.Template, data interface{}) {
	var b bytes.Buffer
	err := tmpl.Execute(&b, data)
	if err != nil {
		glog.Error(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	b.WriteTo(w)
}
