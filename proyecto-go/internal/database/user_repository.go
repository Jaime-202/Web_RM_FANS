package database

import (
	"encoding/json"
	"os"
	"proyecto-go/internal/models"
	"sync"
)

type JSONUserRepository struct {
	filePath string
	mu       sync.RWMutex
}

func NewJSONUserRepository(filePath string) *JSONUserRepository {
	return &JSONUserRepository{
		filePath: filePath,
	}
}

func (repo *JSONUserRepository) readAll() ([]*models.User, error) {
	file, err := os.Open(repo.filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return []*models.User{}, nil
		}
		return nil, err
	}
	defer file.Close()

	var users []*models.User
	if err := json.NewDecoder(file).Decode(&users); err != nil {
		// Archivo vacío
		return []*models.User{}, nil
	}
	return users, nil
}

func (repo *JSONUserRepository) writeAll(users []*models.User) error {
	data, err := json.MarshalIndent(users, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(repo.filePath, data, 0644)
}

func (repo *JSONUserRepository) Save(user *models.User) error {
	repo.mu.Lock()
	defer repo.mu.Unlock()

	users, err := repo.readAll()
	if err != nil {
		return err
	}
	users = append(users, user)
	return repo.writeAll(users)
}

func (repo *JSONUserRepository) FindByUsername(username string) (*models.User, error) {
	repo.mu.RLock()
	defer repo.mu.RUnlock()

	users, err := repo.readAll()
	if err != nil {
		return nil, err
	}
	for _, u := range users {
		if u.Username == username {
			return u, nil
		}
	}
	return nil, nil // Not found
}

func (repo *JSONUserRepository) FindByID(id string) (*models.User, error) {
	repo.mu.RLock()
	defer repo.mu.RUnlock()

	users, err := repo.readAll()
	if err != nil {
		return nil, err
	}
	for _, u := range users {
		if u.ID == id {
			return u, nil
		}
	}
	return nil, nil // Not found
}
