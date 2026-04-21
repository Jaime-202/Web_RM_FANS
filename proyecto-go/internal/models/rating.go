package models

import "time"

// Rating representa la puntuación que un usuario le da a una noticia.
type Rating struct {
	ID        string    `json:"id"`
	UserID    string    `json:"user_id"`
	NewsID    string    `json:"news_id"`
	Score     int       `json:"score"` // 1 a 5 estrellas
	CreatedAt time.Time `json:"created_at"`
}

// RatingRepository define la interfaz para persistir puntuaciones.
type RatingRepository interface {
	Save(rating *Rating) error
	FindByNewsID(newsID string) ([]*Rating, error)
	FindByUserAndNews(userID, newsID string) (*Rating, error)
}
