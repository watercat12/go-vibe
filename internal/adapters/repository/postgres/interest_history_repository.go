package postgres

import (
	"context"
	"e-wallet/internal/domain/interest_history"
	"e-wallet/internal/ports"
	"time"

	"gorm.io/gorm"
)

const (
	InterestHistoryTableName = "interest_history"
)

type interestHistoryRepository struct {
	db *gorm.DB
}

func NewInterestHistoryRepository(db *gorm.DB) ports.InterestHistoryRepository {
	return &interestHistoryRepository{db: db}
}

func (r *interestHistoryRepository) Create(ctx context.Context, ih *interest_history.InterestHistory) (*interest_history.InterestHistory, error) {
	schema := InterestHistory{
		ID:             ih.ID,
		AccountID:      ih.AccountID,
		Date:           ih.Date,
		InterestAmount: ih.InterestAmount,
	}

	if err := r.db.WithContext(ctx).Table(InterestHistoryTableName).Create(&schema).Error; err != nil {
		return nil, err
	}

	return schema.ToDomain(), nil
}

type InterestHistory struct {
	ID             string
	AccountID      string
	Date           time.Time
	InterestAmount float64
	CreatedAt      time.Time `gorm:"autoCreateTime"`
}

func (i *InterestHistory) ToDomain() *interest_history.InterestHistory {
	return &interest_history.InterestHistory{
		ID:             i.ID,
		AccountID:      i.AccountID,
		Date:           i.Date,
		InterestAmount: i.InterestAmount,
		CreatedAt:      i.CreatedAt,
	}
}