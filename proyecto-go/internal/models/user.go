package models

import "time"

// User representa a un usuario del sistema (Admin o Normal).
type User struct {
	ID        string    `json:"id"`
	Username  string    `json:"username"`
	Password  string    `json:"password"` // En un entorno real debería estar hasheada
	Role      string    `json:"role"`     // "admin" o "user"
	CreatedAt time.Time `json:"created_at"`
}

// UserRepository define la interfaz para persistir usuarios.
type UserRepository interface {
	Save(user *User) error
	FindByUsername(username string) (*User, error)
	FindByID(id string) (*User, error)
}
