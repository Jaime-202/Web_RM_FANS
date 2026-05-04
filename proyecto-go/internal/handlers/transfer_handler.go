package handlers

import (
	"encoding/json"
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

func (h *TransferHandler) ServeAPI(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	role, _ := r.Context().Value(middleware.RoleKey).(string)
	isAdmin := role == "admin"

	switch r.Method {
	case http.MethodGet:
		transfers, _ := h.transferRepo.FindAll()
		json.NewEncoder(w).Encode(transfers)

	case http.MethodPost:
		if !isAdmin {
			http.Error(w, `{"error": "Forbidden"}`, http.StatusForbidden)
			return
		}
		var transfer models.Transfer
		if err := json.NewDecoder(r.Body).Decode(&transfer); err != nil {
			http.Error(w, `{"error": "Bad Request"}`, http.StatusBadRequest)
			return
		}
		transfer.ID = uuid.New().String()
		transfer.CreatedAt = time.Now()
		h.transferRepo.Save(&transfer)
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(transfer)

	case http.MethodPut:
		if !isAdmin {
			http.Error(w, `{"error": "Forbidden"}`, http.StatusForbidden)
			return
		}
		var transfer models.Transfer
		if err := json.NewDecoder(r.Body).Decode(&transfer); err != nil {
			http.Error(w, `{"error": "Bad Request"}`, http.StatusBadRequest)
			return
		}
		existing, _ := h.transferRepo.FindByID(transfer.ID)
		if existing == nil {
			http.Error(w, `{"error": "Not Found"}`, http.StatusNotFound)
			return
		}
		transfer.CreatedAt = existing.CreatedAt
		h.transferRepo.Update(&transfer)
		json.NewEncoder(w).Encode(transfer)

	case http.MethodDelete:
		if !isAdmin {
			http.Error(w, `{"error": "Forbidden"}`, http.StatusForbidden)
			return
		}
		id := r.URL.Query().Get("id")
		if id == "" {
			http.Error(w, `{"error": "Missing ID"}`, http.StatusBadRequest)
			return
		}
		existing, _ := h.transferRepo.FindByID(id)
		if existing == nil {
			http.Error(w, `{"error": "Not Found"}`, http.StatusNotFound)
			return
		}
		h.transferRepo.Delete(id)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"success": true}`))

	default:
		http.Error(w, `{"error": "Method Not Allowed"}`, http.StatusMethodNotAllowed)
	}
}
