package handlers

import (
	"crypto/sha256"
	"encoding/hex"
	"html/template"
	"net/http"
	"proyecto-go/internal/models"
	"proyecto-go/internal/session"
	"time"
	"github.com/google/uuid"
)

type AuthHandler struct {
	userRepo models.UserRepository
	tmpl     *template.Template
}

func NewAuthHandler(repo models.UserRepository, t *template.Template) *AuthHandler {
	return &AuthHandler{
		userRepo: repo,
		tmpl:     t,
	}
}

func hashPassword(password string) string {
	hash := sha256.Sum256([]byte(password))
	return hex.EncodeToString(hash[:])
}

func (h *AuthHandler) ServeLogin(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		h.tmpl.ExecuteTemplate(w, "login.html", nil)
		return
	}

	if r.Method == http.MethodPost {
		username := r.FormValue("username")
		password := r.FormValue("password")

		user, err := h.userRepo.FindByUsername(username)
		if err != nil || user == nil || user.Password != hashPassword(password) {
			h.tmpl.ExecuteTemplate(w, "login.html", map[string]string{"Error": "Credenciales inválidas"})
			return
		}

		// Login exitoso, crear sesión
		token, expires := session.GlobalStore.CreateSession(user.ID, user.Role)
		http.SetCookie(w, &http.Cookie{
			Name:     "session_token",
			Value:    token,
			Expires:  expires,
			Path:     "/",
			HttpOnly: true,
		})

		http.Redirect(w, r, "/noticias", http.StatusSeeOther)
	}
}

func (h *AuthHandler) ServeRegister(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		h.tmpl.ExecuteTemplate(w, "register.html", nil)
		return
	}

	if r.Method == http.MethodPost {
		username := r.FormValue("username")
		password := r.FormValue("password")

		// Comprobar si ya existe
		existingUser, _ := h.userRepo.FindByUsername(username)
		if existingUser != nil {
			h.tmpl.ExecuteTemplate(w, "register.html", map[string]string{"Error": "El usuario ya existe"})
			return
		}

		user := &models.User{
			ID:        uuid.New().String(),
			Username:  username,
			Password:  hashPassword(password),
			Role:      "user", // Por defecto todos son usuarios normales
			CreatedAt: time.Now(),
		}

		err := h.userRepo.Save(user)
		if err != nil {
			http.Error(w, "Error al registrar", http.StatusInternalServerError)
			return
		}

		// Auto-login después de registro
		token, expires := session.GlobalStore.CreateSession(user.ID, user.Role)
		http.SetCookie(w, &http.Cookie{
			Name:     "session_token",
			Value:    token,
			Expires:  expires,
			Path:     "/",
			HttpOnly: true,
		})

		http.Redirect(w, r, "/noticias", http.StatusSeeOther)
	}
}

func (h *AuthHandler) ServeLogout(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("session_token")
	if err == nil {
		session.GlobalStore.DeleteSession(cookie.Value)
	}
	http.SetCookie(w, &http.Cookie{
		Name:     "session_token",
		Value:    "",
		Expires:  time.Now().Add(-1 * time.Hour),
		Path:     "/",
		HttpOnly: true,
	})
	http.Redirect(w, r, "/login", http.StatusSeeOther)
}
