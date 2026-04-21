package database

import (
	"encoding/json"
	"errors"
	"os"
	"proyecto-go/internal/models"
	"sync"
)

type JSONPlayerRepository struct {
	filePath string
	mu       sync.RWMutex
}

func NewJSONPlayerRepository(filePath string) *JSONPlayerRepository {
	return &JSONPlayerRepository{
		filePath: filePath,
	}
}

func (repo *JSONPlayerRepository) readAll() ([]*models.Player, error) {
	file, err := os.Open(repo.filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return []*models.Player{}, nil
		}
		return nil, err
	}
	defer file.Close()

	var items []*models.Player
	if err := json.NewDecoder(file).Decode(&items); err != nil {
		return []*models.Player{}, nil
	}
	return items, nil
}

func (repo *JSONPlayerRepository) writeAll(items []*models.Player) error {
	data, err := json.MarshalIndent(items, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(repo.filePath, data, 0644)
}

func (repo *JSONPlayerRepository) Save(n *models.Player) error {
	repo.mu.Lock()
	defer repo.mu.Unlock()

	items, err := repo.readAll()
	if err != nil {
		return err
	}
	items = append(items, n)
	return repo.writeAll(items)
}

func (repo *JSONPlayerRepository) FindAll() ([]*models.Player, error) {
	repo.mu.RLock()
	defer repo.mu.RUnlock()
	return repo.readAll()
}

func (repo *JSONPlayerRepository) FindByID(id string) (*models.Player, error) {
	repo.mu.RLock()
	defer repo.mu.RUnlock()

	items, err := repo.readAll()
	if err != nil {
		return nil, err
	}
	for _, n := range items {
		if n.ID == id {
			return n, nil
		}
	}
	return nil, nil // Not found
}

func (repo *JSONPlayerRepository) Delete(id string) error {
	repo.mu.Lock()
	defer repo.mu.Unlock()

	items, err := repo.readAll()
	if err != nil {
		return err
	}

	for i, n := range items {
		if n.ID == id {
			items = append(items[:i], items[i+1:]...)
			return repo.writeAll(items)
		}
	}
	return errors.New("not found")
}

func (repo *JSONPlayerRepository) Update(updatedItem *models.Player) error {
	repo.mu.Lock()
	defer repo.mu.Unlock()

	items, err := repo.readAll()
	if err != nil {
		return err
	}

	for i, n := range items {
		if n.ID == updatedItem.ID {
			items[i] = updatedItem
			return repo.writeAll(items)
		}
	}
	return errors.New("not found")
}
