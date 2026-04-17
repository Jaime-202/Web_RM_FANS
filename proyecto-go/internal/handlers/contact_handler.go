package handlers

import (
	"html/template"
	"log"
	"net/http"
	"path/filepath"
	"proyecto-go/internal/services"
)

type ContactHandler struct {
	service *services.ContactService
}

func NewContactHandler(service *services.ContactService) *ContactHandler {
	return &ContactHandler{
		service: service,
	}
}

// ServeForm pinta el formulario HTML cuando se hace un GET a /contacto
func (h *ContactHandler) ServeForm(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	tmplPath := filepath.Join("web", "templates", "contact.html")
	tmpl, err := template.ParseFiles(tmplPath)
	if err != nil {
		http.Error(w, "Error cargando la vista", http.StatusInternalServerError)
		log.Printf("Error al cargar la plantilla: %v", err)
		return
	}

	tmpl.Execute(w, nil)
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

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("¡Gracias por contactar con nuestra tienda de mascotas! Hemos recibido tu mensaje."))
}
