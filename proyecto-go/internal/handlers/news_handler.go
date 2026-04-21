package handlers

import (
	"fmt"
	"html/template"
	"net/http"
	"proyecto-go/internal/middleware"
	"proyecto-go/internal/models"
	"strconv"
	"time"

	"github.com/google/uuid"
)

type NewsHandler struct {
	newsRepo   models.NewsRepository
	ratingRepo models.RatingRepository
	tmpl       *template.Template
}

func NewNewsHandler(nRepo models.NewsRepository, rRepo models.RatingRepository, t *template.Template) *NewsHandler {
	return &NewsHandler{
		newsRepo:   nRepo,
		ratingRepo: rRepo,
		tmpl:       t,
	}
}

func (h *NewsHandler) ServeList(w http.ResponseWriter, r *http.Request) {
	newsList, err := h.newsRepo.FindAll()
	if err != nil {
		http.Error(w, "Error al cargar noticias", http.StatusInternalServerError)
		return
	}

	role, _ := r.Context().Value(middleware.RoleKey).(string)
	isAdmin := role == "admin"

	data := map[string]interface{}{
		"News":    newsList,
		"IsAdmin": isAdmin,
		"IsLoggedIn": r.Context().Value(middleware.UserIDKey) != nil,
	}

	h.tmpl.ExecuteTemplate(w, "news_list.html", data)
}

func (h *NewsHandler) ServeDetail(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		http.Redirect(w, r, "/noticias", http.StatusSeeOther)
		return
	}

	news, err := h.newsRepo.FindByID(id)
	if err != nil || news == nil {
		http.Error(w, "Noticia no encontrada", http.StatusNotFound)
		return
	}

	userID, _ := r.Context().Value(middleware.UserIDKey).(string)
	
	// Buscar si este usuario ya ha votado esta noticia
	var userRating *models.Rating
	if userID != "" {
		userRating, _ = h.ratingRepo.FindByUserAndNews(userID, id)
	}

	data := map[string]interface{}{
		"News":       news,
		"IsLoggedIn": userID != "",
		"UserRating": userRating,
	}

	h.tmpl.ExecuteTemplate(w, "news_detail.html", data)
}

// Admin: Crear noticia
func (h *NewsHandler) ServeCreateForm(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		data := map[string]interface{}{
			"Action": "Crear",
			"News":   &models.News{},
			"IsAdmin": true,
			"IsLoggedIn": true,
		}
		h.tmpl.ExecuteTemplate(w, "news_form.html", data)
		return
	}

	if r.Method == http.MethodPost {
		news := &models.News{
			ID:        uuid.New().String(),
			Title:     r.FormValue("title"),
			Summary:   r.FormValue("summary"),
			Content:   r.FormValue("content"),
			ImageURL:  r.FormValue("image_url"),
			CreatedAt: time.Now(),
		}
		
		h.newsRepo.Save(news)
		http.Redirect(w, r, "/noticias", http.StatusSeeOther)
	}
}

// Admin: Editar noticia
func (h *NewsHandler) ServeEditForm(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	news, err := h.newsRepo.FindByID(id)
	if err != nil || news == nil {
		http.Error(w, "Noticia no encontrada", http.StatusNotFound)
		return
	}

	if r.Method == http.MethodGet {
		data := map[string]interface{}{
			"Action": "Editar",
			"News":   news,
			"IsAdmin": true,
			"IsLoggedIn": true,
		}
		h.tmpl.ExecuteTemplate(w, "news_form.html", data)
		return
	}

	if r.Method == http.MethodPost {
		news.Title = r.FormValue("title")
		news.Summary = r.FormValue("summary")
		news.Content = r.FormValue("content")
		news.ImageURL = r.FormValue("image_url")
		h.newsRepo.Update(news)
		http.Redirect(w, r, "/noticias", http.StatusSeeOther)
	}
}

// Admin: Borrar noticia
func (h *NewsHandler) ServeDelete(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id != "" {
		h.newsRepo.Delete(id)
	}
	http.Redirect(w, r, "/noticias", http.StatusSeeOther)
}

// Usuario: Enviar puntuación (Rating)
func (h *NewsHandler) ProcessRating(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	userID := r.Context().Value(middleware.UserIDKey).(string)
	newsID := r.FormValue("news_id")
	scoreStr := r.FormValue("score")
	
	score, err := strconv.Atoi(scoreStr)
	if err != nil || score < 1 || score > 5 {
		http.Error(w, "Puntuación inválida", http.StatusBadRequest)
		return
	}

	rating := &models.Rating{
		ID:        uuid.New().String(),
		UserID:    userID,
		NewsID:    newsID,
		Score:     score,
		CreatedAt: time.Now(),
	}

	// Guardar el rating
	err = h.ratingRepo.Save(rating)
	if err == nil {
		// Recalcular la media
		h.updateNewsAverageRating(newsID)
	}

	// Redirigir de vuelta a la noticia
	http.Redirect(w, r, fmt.Sprintf("/noticia?id=%s", newsID), http.StatusSeeOther)
}

func (h *NewsHandler) updateNewsAverageRating(newsID string) {
	ratings, err := h.ratingRepo.FindByNewsID(newsID)
	if err != nil || len(ratings) == 0 {
		return
	}

	var sum int
	for _, r := range ratings {
		sum += r.Score
	}
	
	avg := float64(sum) / float64(len(ratings))

	news, err := h.newsRepo.FindByID(newsID)
	if err == nil && news != nil {
		news.AverageRating = avg
		news.RatingCount = len(ratings)
		h.newsRepo.Update(news)
	}
}
