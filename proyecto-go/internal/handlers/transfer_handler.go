package handlers

import (
	"html/template"
	"net/http"
	"proyecto-go/internal/middleware"
	"proyecto-go/internal/models"
	"time"

	"github.com/google/uuid"
)

type TransferHandler struct {
	transferRepo models.TransferRepository
	tmpl         *template.Template
}

func NewTransferHandler(tr models.TransferRepository, t *template.Template) *TransferHandler {
	return &TransferHandler{
		transferRepo: tr,
		tmpl:         t,
	}
}

func (h *TransferHandler) ServeList(w http.ResponseWriter, r *http.Request) {
	transfers, _ := h.transferRepo.FindAll()
	role, _ := r.Context().Value(middleware.RoleKey).(string)

	data := map[string]interface{}{
		"IsLoggedIn": r.Context().Value(middleware.UserIDKey) != nil,
		"IsAdmin":    role == "admin",
		"Transfers":  transfers,
	}

	h.tmpl.ExecuteTemplate(w, "fichajes.html", data)
}

func (h *TransferHandler) ServeCreateForm(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		data := map[string]interface{}{
			"Action":     "Crear",
			"Transfer":   &models.Transfer{},
			"IsAdmin":    true,
			"IsLoggedIn": true,
		}
		h.tmpl.ExecuteTemplate(w, "transfer_form.html", data)
		return
	}

	if r.Method == http.MethodPost {
		transfer := &models.Transfer{
			ID:          uuid.New().String(),
			PlayerName:  r.FormValue("player_name"),
			FromTeam:    r.FormValue("from_team"),
			ToTeam:      r.FormValue("to_team"),
			Status:      r.FormValue("status"),
			Description: r.FormValue("description"),
			ImageURL:    r.FormValue("image_url"),
			CreatedAt:   time.Now(),
		}
		h.transferRepo.Save(transfer)
		http.Redirect(w, r, "/fichajes", http.StatusSeeOther)
	}
}

func (h *TransferHandler) ServeEditForm(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	transfer, _ := h.transferRepo.FindByID(id)
	if transfer == nil {
		http.NotFound(w, r)
		return
	}

	if r.Method == http.MethodGet {
		data := map[string]interface{}{
			"Action":     "Editar",
			"Transfer":   transfer,
			"IsAdmin":    true,
			"IsLoggedIn": true,
		}
		h.tmpl.ExecuteTemplate(w, "transfer_form.html", data)
		return
	}

	if r.Method == http.MethodPost {
		transfer.PlayerName = r.FormValue("player_name")
		transfer.FromTeam = r.FormValue("from_team")
		transfer.ToTeam = r.FormValue("to_team")
		transfer.Status = r.FormValue("status")
		transfer.Description = r.FormValue("description")
		transfer.ImageURL = r.FormValue("image_url")
		h.transferRepo.Update(transfer)
		http.Redirect(w, r, "/fichajes", http.StatusSeeOther)
	}
}

func (h *TransferHandler) ServeDelete(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	h.transferRepo.Delete(id)
	http.Redirect(w, r, "/fichajes", http.StatusSeeOther)
}
