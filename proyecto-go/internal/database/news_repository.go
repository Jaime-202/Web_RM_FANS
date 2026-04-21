package database

import (
	"encoding/json"
	"errors"
	"os"
	"proyecto-go/internal/models"
	"sync"
)

type JSONNewsRepository struct {
	filePath string
	mu       sync.RWMutex
}

func NewJSONNewsRepository(filePath string) *JSONNewsRepository {
	return &JSONNewsRepository{
		filePath: filePath,
	}
}

func (repo *JSONNewsRepository) readAll() ([]*models.News, error) {
	file, err := os.Open(repo.filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return []*models.News{}, nil
		}
		return nil, err
	}
	defer file.Close()

	var news []*models.News
	if err := json.NewDecoder(file).Decode(&news); err != nil {
		return []*models.News{}, nil
	}
	return news, nil
}

func (repo *JSONNewsRepository) writeAll(news []*models.News) error {
	data, err := json.MarshalIndent(news, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(repo.filePath, data, 0644)
}

func (repo *JSONNewsRepository) Save(n *models.News) error {
	repo.mu.Lock()
	defer repo.mu.Unlock()

	news, err := repo.readAll()
	if err != nil {
		return err
	}
	news = append(news, n)
	return repo.writeAll(news)
}

func (repo *JSONNewsRepository) FindAll() ([]*models.News, error) {
	repo.mu.RLock()
	defer repo.mu.RUnlock()
	return repo.readAll()
}

func (repo *JSONNewsRepository) FindByID(id string) (*models.News, error) {
	repo.mu.RLock()
	defer repo.mu.RUnlock()

	news, err := repo.readAll()
	if err != nil {
		return nil, err
	}
	for _, n := range news {
		if n.ID == id {
			return n, nil
		}
	}
	return nil, nil // Not found
}

func (repo *JSONNewsRepository) Delete(id string) error {
	repo.mu.Lock()
	defer repo.mu.Unlock()

	news, err := repo.readAll()
	if err != nil {
		return err
	}

	for i, n := range news {
		if n.ID == id {
			news = append(news[:i], news[i+1:]...)
			return repo.writeAll(news)
		}
	}
	return errors.New("news not found")
}

func (repo *JSONNewsRepository) Update(updatedNews *models.News) error {
	repo.mu.Lock()
	defer repo.mu.Unlock()

	news, err := repo.readAll()
	if err != nil {
		return err
	}

	for i, n := range news {
		if n.ID == updatedNews.ID {
			news[i] = updatedNews
			return repo.writeAll(news)
		}
	}
	return errors.New("news not found")
}
