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

type PlayerHandler struct {
	playerRepo models.PlayerRepository
	tmpl       *template.Template
}

func NewPlayerHandler(pr models.PlayerRepository, t *template.Template) *PlayerHandler {
	return &PlayerHandler{
		playerRepo: pr,
		tmpl:       t,
	}
}

func (h *PlayerHandler) ServeList(w http.ResponseWriter, r *http.Request) {
	players, _ := h.playerRepo.FindAll()
	role, _ := r.Context().Value(middleware.RoleKey).(string)

	// Agrupar jugadores por posición
	groups := map[string][]*models.Player{
		"Portero":        {},
		"Defensa":        {},
		"Centrocampista": {},
		"Delantero":      {},
	}
	for _, p := range players {
		if _, ok := groups[p.Position]; ok {
			groups[p.Position] = append(groups[p.Position], p)
		}
	}

	data := map[string]interface{}{
		"IsLoggedIn":      r.Context().Value(middleware.UserIDKey) != nil,
		"IsAdmin":         role == "admin",
		"Players":         players,
		"Porteros":        groups["Portero"],
		"Defensas":        groups["Defensa"],
		"Centrocampistas": groups["Centrocampista"],
		"Delanteros":      groups["Delantero"],
	}

	h.tmpl.ExecuteTemplate(w, "plantilla.html", data)
}

func (h *PlayerHandler) ServeCreateForm(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		data := map[string]interface{}{
			"Action":     "Crear",
			"Player":     &models.Player{},
			"IsAdmin":    true,
			"IsLoggedIn": true,
		}
		h.tmpl.ExecuteTemplate(w, "player_form.html", data)
		return
	}

	if r.Method == http.MethodPost {
		player := &models.Player{
			ID:        uuid.New().String(),
			Name:      r.FormValue("name"),
			Position:  r.FormValue("position"),
			ImageURL:  r.FormValue("image_url"),
			CreatedAt: time.Now(),
		}
		h.playerRepo.Save(player)
		http.Redirect(w, r, "/plantilla", http.StatusSeeOther)
	}
}

func (h *PlayerHandler) ServeEditForm(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	player, _ := h.playerRepo.FindByID(id)
	if player == nil {
		http.NotFound(w, r)
		return
	}

	if r.Method == http.MethodGet {
		data := map[string]interface{}{
			"Action":     "Editar",
			"Player":     player,
			"IsAdmin":    true,
			"IsLoggedIn": true,
		}
		h.tmpl.ExecuteTemplate(w, "player_form.html", data)
		return
	}

	if r.Method == http.MethodPost {
		player.Name = r.FormValue("name")
		player.Position = r.FormValue("position")
		player.ImageURL = r.FormValue("image_url")
		h.playerRepo.Update(player)
		http.Redirect(w, r, "/plantilla", http.StatusSeeOther)
	}
}

func (h *PlayerHandler) ServeDelete(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	h.playerRepo.Delete(id)
	http.Redirect(w, r, "/plantilla", http.StatusSeeOther)
}

func (h *PlayerHandler) ServeAPI(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	role, _ := r.Context().Value(middleware.RoleKey).(string)
	isAdmin := role == "admin"

	switch r.Method {
	case http.MethodGet:
		players, _ := h.playerRepo.FindAll()
		json.NewEncoder(w).Encode(players)

	case http.MethodPost:
		if !isAdmin {
			http.Error(w, `{"error": "Forbidden"}`, http.StatusForbidden)
			return
		}
		var player models.Player
		if err := json.NewDecoder(r.Body).Decode(&player); err != nil {
			http.Error(w, `{"error": "Bad Request"}`, http.StatusBadRequest)
			return
		}
		player.ID = uuid.New().String()
		player.CreatedAt = time.Now()
		h.playerRepo.Save(&player)
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(player)

	case http.MethodPut:
		if !isAdmin {
			http.Error(w, `{"error": "Forbidden"}`, http.StatusForbidden)
			return
		}
		var player models.Player
		if err := json.NewDecoder(r.Body).Decode(&player); err != nil {
			http.Error(w, `{"error": "Bad Request"}`, http.StatusBadRequest)
			return
		}
		existing, _ := h.playerRepo.FindByID(player.ID)
		if existing == nil {
			http.Error(w, `{"error": "Not Found"}`, http.StatusNotFound)
			return
		}
		player.CreatedAt = existing.CreatedAt
		h.playerRepo.Update(&player)
		json.NewEncoder(w).Encode(player)

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
		existing, _ := h.playerRepo.FindByID(id)
		if existing == nil {
			http.Error(w, `{"error": "Not Found"}`, http.StatusNotFound)
			return
		}
		h.playerRepo.Delete(id)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"success": true}`))

	default:
		http.Error(w, `{"error": "Method Not Allowed"}`, http.StatusMethodNotAllowed)
	}
}
