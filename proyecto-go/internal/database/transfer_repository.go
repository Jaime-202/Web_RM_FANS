package database

import (
	"encoding/json"
	"errors"
	"os"
	"proyecto-go/internal/models"
	"sync"
)

type JSONTransferRepository struct {
	filePath string
	mu       sync.RWMutex
}

func NewJSONTransferRepository(filePath string) *JSONTransferRepository {
	return &JSONTransferRepository{
		filePath: filePath,
	}
}

func (repo *JSONTransferRepository) readAll() ([]*models.Transfer, error) {
	file, err := os.Open(repo.filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return []*models.Transfer{}, nil
		}
		return nil, err
	}
	defer file.Close()

	var items []*models.Transfer
	if err := json.NewDecoder(file).Decode(&items); err != nil {
		return []*models.Transfer{}, nil
	}
	return items, nil
}

func (repo *JSONTransferRepository) writeAll(items []*models.Transfer) error {
	data, err := json.MarshalIndent(items, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(repo.filePath, data, 0644)
}

func (repo *JSONTransferRepository) Save(n *models.Transfer) error {
	repo.mu.Lock()
	defer repo.mu.Unlock()

	items, err := repo.readAll()
	if err != nil {
		return err
	}
	items = append(items, n)
	return repo.writeAll(items)
}

func (repo *JSONTransferRepository) FindAll() ([]*models.Transfer, error) {
	repo.mu.RLock()
	defer repo.mu.RUnlock()
	return repo.readAll()
}

func (repo *JSONTransferRepository) FindByID(id string) (*models.Transfer, error) {
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

func (repo *JSONTransferRepository) Delete(id string) error {
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

func (repo *JSONTransferRepository) Update(updatedItem *models.Transfer) error {
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
