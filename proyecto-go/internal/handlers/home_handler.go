package handlers

import (
	"html/template"
	"net/http"
	"proyecto-go/internal/middleware"
	"proyecto-go/internal/models"
)

type HomeHandler struct {
	newsRepo     models.NewsRepository
	playerRepo   models.PlayerRepository
	transferRepo models.TransferRepository
	tmpl         *template.Template
}

func NewHomeHandler(nr models.NewsRepository, pr models.PlayerRepository, tr models.TransferRepository, t *template.Template) *HomeHandler {
	return &HomeHandler{
		newsRepo:     nr,
		playerRepo:   pr,
		transferRepo: tr,
		tmpl:         t,
	}
}

func (h *HomeHandler) ServeHome(w http.ResponseWriter, r *http.Request) {
	news, _ := h.newsRepo.FindAll()
	allPlayers, _ := h.playerRepo.FindAll()
	allTransfers, _ := h.transferRepo.FindAll()

	role, _ := r.Context().Value(middleware.RoleKey).(string)

	// Noticia destacada (la primera)
	var featuredNews *models.News
	if len(news) > 0 {
		featuredNews = news[0]
	}

	// Limitar la plantilla a 4 jugadores en el inicio
	previewPlayers := allPlayers
	if len(previewPlayers) > 4 {
		previewPlayers = previewPlayers[:4]
	}

	// Limitar los fichajes a 3 en el inicio
	previewTransfers := allTransfers
	if len(previewTransfers) > 3 {
		previewTransfers = previewTransfers[:3]
	}

	data := map[string]interface{}{
		"IsLoggedIn":   r.Context().Value(middleware.UserIDKey) != nil,
		"IsAdmin":      role == "admin",
		"News":         news,
		"FeaturedNews": featuredNews,
		"Players":      previewPlayers,
		"Transfers":    previewTransfers,
	}

	h.tmpl.ExecuteTemplate(w, "index.html", data)
}
