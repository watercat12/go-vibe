package postgres

import (
	"context"
	"e-wallet/internal/domain/transaction"
	"e-wallet/internal/ports"
	"time"

	"gorm.io/gorm"
)

const (
	TransactionsTableName = "transactions"
)

type transactionRepository struct {
	db *gorm.DB
}

func NewTransactionRepository(db *gorm.DB) ports.TransactionRepository {
	return &transactionRepository{db: db}
}

func (r *transactionRepository) Create(ctx context.Context, tx *transaction.Transaction) (*transaction.Transaction, error) {
	schema := Transaction{
		ID:               tx.ID,
		AccountID:        tx.AccountID,
		TransactionType:  tx.TransactionType,
		Amount:           tx.Amount,
		Status:           tx.Status,
		BalanceAfter:     tx.BalanceAfter,
		RelatedAccountID: tx.RelatedAccountID,
	}

	if err := r.db.WithContext(ctx).Table(TransactionsTableName).Create(&schema).Error; err != nil {
		return nil, err
	}

	return schema.ToDomain(), nil
}

type Transaction struct {
	ID               string
	AccountID        string
	TransactionType  string
	Amount           float64
	Status           string
	BalanceAfter     float64
	RelatedAccountID *string
	CreatedAt        time.Time `gorm:"autoCreateTime"`
}

func (t *Transaction) ToDomain() *transaction.Transaction {
	return &transaction.Transaction{
		ID:               t.ID,
		AccountID:        t.AccountID,
		TransactionType:  t.TransactionType,
		Amount:           t.Amount,
		Status:           t.Status,
		BalanceAfter:     t.BalanceAfter,
		RelatedAccountID: t.RelatedAccountID,
		CreatedAt:        t.CreatedAt,
	}
}