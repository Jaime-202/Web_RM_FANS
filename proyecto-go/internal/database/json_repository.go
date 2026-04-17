package database

import (
	"encoding/json"
	"os"
	"proyecto-go/internal/models"
	"sync"
)

// JSONContactRepository implementa models.ContactRepository y guarda en JSONL.
type JSONContactRepository struct {
	filePath string
	mu       sync.Mutex // Previene problemas de concurrencia al escribir en archivo
}

func NewJSONContactRepository(filePath string) *JSONContactRepository {
	return &JSONContactRepository{
		filePath: filePath,
	}
}

// Save guarda un nuevo contacto en el archivo JSONL (JSON Lines).
func (repo *JSONContactRepository) Save(contact *models.Contact) error {
	repo.mu.Lock()
	defer repo.mu.Unlock()

	file, err := os.OpenFile(repo.filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	return encoder.Encode(contact)
}
