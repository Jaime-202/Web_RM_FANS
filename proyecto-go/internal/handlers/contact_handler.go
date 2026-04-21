package handlers

import (
	"html/template"
	"net/http"
	"proyecto-go/internal/middleware"
	"proyecto-go/internal/services"
)

type ContactHandler struct {
	service *services.ContactService
	tmpl    *template.Template
}

func NewContactHandler(service *services.ContactService, t *template.Template) *ContactHandler {
	return &ContactHandler{
		service: service,
		tmpl:    t,
	}
}

// ServeForm pinta el formulario HTML cuando se hace un GET a /contacto
func (h *ContactHandler) ServeForm(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	role, _ := r.Context().Value(middleware.RoleKey).(string)
	data := map[string]interface{}{
		"IsLoggedIn": r.Context().Value(middleware.UserIDKey) != nil,
		"IsAdmin":    role == "admin",
	}

	h.tmpl.ExecuteTemplate(w, "contact.html", data)
}

// ProcessForm atiende el POST proveniente del formulario HTML
func (h *ContactHandler) ProcessForm(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Error procesando petición", http.StatusBadRequest)
		return
	}

	name := r.FormValue("name")
	email := r.FormValue("email")
	message := r.FormValue("message")

	err = h.service.ProcessContactCreation(name, email, message)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	http.Redirect(w, r, "/?contacto=ok", http.StatusSeeOther)
}
