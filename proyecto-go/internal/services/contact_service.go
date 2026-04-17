package services

import (
	"errors"
	"fmt"
	"proyecto-go/internal/models"
	"time"
)

type ContactService struct {
	repo models.ContactRepository
}

func NewContactService(repo models.ContactRepository) *ContactService {
	return &ContactService{
		repo: repo,
	}
}

// ProcessContactCreation valida los datos del formulario y los persiste.
func (s *ContactService) ProcessContactCreation(name, email, message string) error {
	if name == "" || email == "" || message == "" {
		return errors.New("todos los campos son obligatorios")
	}

	contact := &models.Contact{
		ID:        fmt.Sprintf("ID-%d", time.Now().UnixNano()), // Simulación de UUID único
		Name:      name,
		Email:     email,
		Message:   message,
		CreatedAt: time.Now(),
	}

	return s.repo.Save(contact)
}
