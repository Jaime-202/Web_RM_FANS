package handlers

import (
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
