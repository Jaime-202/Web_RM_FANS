package models

import "time"

// Contact representa la información de un contacto recibido en la tienda de mascotas.
type Contact struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Message   string    `json:"message"`
	CreatedAt time.Time `json:"created_at"`
}

// ContactRepository define la interfaz para la persistencia, permitiendo
// intercambiar la base de datos (JSON, MySQL, PostgreSQL) fácilmente en el futuro.
type ContactRepository interface {
	Save(contact *Contact) error
}
