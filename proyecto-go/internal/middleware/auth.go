package middleware

import (
	"context"
	"net/http"
	"proyecto-go/internal/session"
)

// Claves para el contexto
type contextKey string

const (
	UserIDKey contextKey = "userID"
	RoleKey   contextKey = "role"
)

// AuthMiddleware verifica si el usuario está logueado leyendo la cookie de sesión.
// Si está logueado, inyecta su ID y Rol en el contexto de la petición.
func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("session_token")
		if err != nil {
			// No hay cookie, continúa sin usuario
			next.ServeHTTP(w, r)
			return
		}

		sess, exists := session.GlobalStore.GetSession(cookie.Value)
		if !exists {
			// Cookie inválida o expirada
			next.ServeHTTP(w, r)
			return
		}

		// Inyectar datos en el contexto
		ctx := context.WithValue(r.Context(), UserIDKey, sess.UserID)
		ctx = context.WithValue(ctx, RoleKey, sess.Role)
		
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// RequireAuth bloquea el acceso si no hay un usuario logueado.
func RequireAuth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID := r.Context().Value(UserIDKey)
		if userID == nil {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}
		next.ServeHTTP(w, r)
	}
}

// RequireAdmin bloquea el acceso si el usuario no tiene el rol "admin".
func RequireAdmin(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		role := r.Context().Value(RoleKey)
		if role == nil || role.(string) != "admin" {
			http.Error(w, "Acceso denegado: Se requieren permisos de administrador", http.StatusForbidden)
			return
		}
		next.ServeHTTP(w, r)
	}
}
