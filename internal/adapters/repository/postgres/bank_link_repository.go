package postgres

import (
	"context"
	"e-wallet/internal/domain/bank_link"
	"e-wallet/internal/ports"

	"gorm.io/gorm"
)

type bankLinkRepository struct {
	db *gorm.DB
}

func NewBankLinkRepository(db *gorm.DB) ports.BankLinkRepository {
	return &bankLinkRepository{db: db}
}

func (r *bankLinkRepository) Create(ctx context.Context, bankLink *bank_link.BankLink) (*bank_link.BankLink, error) {
	if err := r.db.Table("bank_links").Create(bankLink).Error; err != nil {
		return nil, err
	}
	return bankLink, nil
}

func (r *bankLinkRepository) GetByUserID(ctx context.Context, userID string) ([]*bank_link.BankLink, error) {
	var bankLinks []*bank_link.BankLink
	if err := r.db.Table("bank_links").Where("user_id = ?", userID).Find(&bankLinks).Error; err != nil {
		return nil, err
	}
	return bankLinks, nil
}

func (r *bankLinkRepository) CountByUserID(ctx context.Context, userID string) (int, error) {
	var count int64
	if err := r.db.Table("bank_links").Where("user_id = ?", userID).Count(&count).Error; err != nil {
		return 0, err
	}
	return int(count), nil
}