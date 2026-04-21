package database

import (
	"encoding/json"
	"os"
	"proyecto-go/internal/models"
	"sync"
)

type JSONRatingRepository struct {
	filePath string
	mu       sync.RWMutex
}

func NewJSONRatingRepository(filePath string) *JSONRatingRepository {
	return &JSONRatingRepository{
		filePath: filePath,
	}
}

func (repo *JSONRatingRepository) readAll() ([]*models.Rating, error) {
	file, err := os.Open(repo.filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return []*models.Rating{}, nil
		}
		return nil, err
	}
	defer file.Close()

	var ratings []*models.Rating
	if err := json.NewDecoder(file).Decode(&ratings); err != nil {
		return []*models.Rating{}, nil
	}
	return ratings, nil
}

func (repo *JSONRatingRepository) writeAll(ratings []*models.Rating) error {
	data, err := json.MarshalIndent(ratings, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(repo.filePath, data, 0644)
}

func (repo *JSONRatingRepository) Save(r *models.Rating) error {
	repo.mu.Lock()
	defer repo.mu.Unlock()

	ratings, err := repo.readAll()
	if err != nil {
		return err
	}
	
	// Si el usuario ya había votado, actualizamos su voto
	found := false
	for i, existing := range ratings {
		if existing.UserID == r.UserID && existing.NewsID == r.NewsID {
			ratings[i].Score = r.Score
			found = true
			break
		}
	}
	
	if !found {
		ratings = append(ratings, r)
	}
	
	return repo.writeAll(ratings)
}

func (repo *JSONRatingRepository) FindByNewsID(newsID string) ([]*models.Rating, error) {
	repo.mu.RLock()
	defer repo.mu.RUnlock()

	ratings, err := repo.readAll()
	if err != nil {
		return nil, err
	}

	var result []*models.Rating
	for _, r := range ratings {
		if r.NewsID == newsID {
			result = append(result, r)
		}
	}
	return result, nil
}

func (repo *JSONRatingRepository) FindByUserAndNews(userID, newsID string) (*models.Rating, error) {
	repo.mu.RLock()
	defer repo.mu.RUnlock()

	ratings, err := repo.readAll()
	if err != nil {
		return nil, err
	}

	for _, r := range ratings {
		if r.UserID == userID && r.NewsID == newsID {
			return r, nil
		}
	}
	return nil, nil // Not found
}
