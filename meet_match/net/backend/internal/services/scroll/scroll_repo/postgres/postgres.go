package postgres

import (
	"test_backend_frontend/internal/models"
	"test_backend_frontend/internal/models/models_da"
	"test_backend_frontend/internal/services/scroll/scroll_repo"
	"time"

	"github.com/google/uuid"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

const MAX_LIMIT = 100

type scrollRepository struct {
	db *gorm.DB
}

func NewScrollRepository(db *gorm.DB) scroll_repo.ScrollRepository {
	return &scrollRepository{db: db}
}

func (s scrollRepository) AddScrollFact(fact *models.FactScrolled) error {
	pgFact := models_da.ToPostgresFactScrolled(fact)
	pgFact.DateTime = time.Now()
	tx := s.db.Create(&pgFact)
	if tx.Error != nil {
		return errors.Wrap(tx.Error, "scroll.repository.AddScrollFact error")
	}
	return nil
}

func (s scrollRepository) GetAllLikedPlaces(session_id uuid.UUID, user_id uint64) ([]uint64, error) {
	var ids []uint64

	tx := s.db.Limit(MAX_LIMIT).
		Model(&models_da.FactScrolled{}).
		Select("place_id").
		Where("session_id = ? AND user_id = ? AND is_liked = true", session_id.String(), user_id).
		Order("datetime desc").
		Find(&ids)

	if tx.Error != nil {
		return nil,
			errors.Wrap(tx.Error, "scroll.repository.GetAllLikedPlaces error")
	}

	return ids, nil
}

func (s scrollRepository) GetAllUsersIdsForSession(session_id uuid.UUID) ([]uint64, error) {
	var ids []uint64

	tx := s.db.Limit(MAX_LIMIT).
		Model(&models_da.FactScrolled{}).
		Select("user_id").
		Where("session_id = ?", session_id.String()).
		Group("user_id").
		Find(&ids)

	if tx.Error != nil {
		return nil,
			errors.Wrap(tx.Error, "scroll.repository.GetAllUsersIdsForSession error")
	}

	return ids, nil
}
