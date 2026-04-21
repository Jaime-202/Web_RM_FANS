package models

import "time"

type Player struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Position  string    `json:"position"` // Portero, Defensa, Centrocampista, Delantero
	ImageURL  string    `json:"image_url"`
	CreatedAt time.Time `json:"created_at"`
}

type PlayerRepository interface {
	Save(player *Player) error
	FindAll() ([]*Player, error)
	FindByID(id string) (*Player, error)
	Delete(id string) error
	Update(player *Player) error
}
