package models

import "time"

type Transfer struct {
	ID          string    `json:"id"`
	PlayerName  string    `json:"player_name"`
	FromTeam    string    `json:"from_team"`
	ToTeam      string    `json:"to_team"`
	Status      string    `json:"status"` // Rumor, Hecho, Contrato verbal, etc
	Description string    `json:"description"`
	ImageURL    string    `json:"image_url"`
	CreatedAt   time.Time `json:"created_at"`
}

type TransferRepository interface {
	Save(transfer *Transfer) error
	FindAll() ([]*Transfer, error)
	FindByID(id string) (*Transfer, error)
	Delete(id string) error
	Update(transfer *Transfer) error
}
