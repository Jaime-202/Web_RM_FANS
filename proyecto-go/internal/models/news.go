package models

import "time"

// News representa una noticia del portal.
type News struct {
	ID            string    `json:"id"`
	Title         string    `json:"title"`
	Summary       string    `json:"summary"`
	Content       string    `json:"content"`
	ImageURL      string    `json:"image_url"`
	AverageRating float64   `json:"average_rating"` // Calculado a partir de los ratings
	RatingCount   int       `json:"rating_count"`   // Número total de votos
	CreatedAt     time.Time `json:"created_at"`
}

// NewsRepository define la interfaz para persistir noticias.
type NewsRepository interface {
	Save(news *News) error
	FindAll() ([]*News, error)
	FindByID(id string) (*News, error)
	Delete(id string) error
	Update(news *News) error
}
