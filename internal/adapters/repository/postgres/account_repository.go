package postgres

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"time"

	"e-wallet/internal/domain/account"
	"e-wallet/internal/ports"
	"e-wallet/pkg"

	"gorm.io/gorm"
)


type accountRepository struct {
	db *gorm.DB
}

type savingsAccountDetailRepository struct {
	db *gorm.DB
}

func NewAccountRepository(db *gorm.DB) ports.AccountRepository {
	return &accountRepository{db: db}
}

func NewSavingsAccountDetailRepository(db *gorm.DB) ports.SavingsAccountDetailRepository {
	return &savingsAccountDetailRepository{db: db}
}

// Account schema
type Account struct {
	ID            string    `gorm:"column:id;primaryKey"`
	UserID        string    `gorm:"column:user_id;not null"`
	AccountNumber string    `gorm:"column:account_number;unique;not null"`
	AccountType   string    `gorm:"column:account_type;not null"`
	Balance       float64   `gorm:"column:balance;default:0"`
	Status        string    `gorm:"column:status;default:ACTIVE"`
	CreatedAt     time.Time `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt     time.Time `gorm:"column:updated_at;autoUpdateTime"`
}

func (a *Account) ToDomain() *account.Account {
	return &account.Account{
		ID:            a.ID,
		UserID:        a.UserID,
		AccountNumber: a.AccountNumber,
		AccountType:   a.AccountType,
		Balance:       a.Balance,
		Status:        a.Status,
		CreatedAt:     a.CreatedAt,
		UpdatedAt:     a.UpdatedAt,
	}
}

// SavingsAccountDetail schema
type SavingsAccountDetail struct {
	AccountID             string     `gorm:"column:account_id;primaryKey"`
	IsFixedTerm           bool       `gorm:"column:is_fixed_term;not null"`
	TermMonths            *int       `gorm:"column:term_months"`
	AnnualInterestRate    float64    `gorm:"column:annual_interest_rate;not null"`
	StartDate             time.Time  `gorm:"column:start_date;not null"`
	MaturityDate          *time.Time `gorm:"column:maturity_date"`
	LastInterestCalcDate  *time.Time `gorm:"column:last_interest_calc_date"`
	CreatedAt             time.Time  `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt             time.Time  `gorm:"column:updated_at;autoUpdateTime"`
}

func (s *SavingsAccountDetail) ToDomain() *account.SavingsAccountDetail {
	return &account.SavingsAccountDetail{
		AccountID:             s.AccountID,
		IsFixedTerm:           s.IsFixedTerm,
		TermMonths:            s.TermMonths,
		AnnualInterestRate:    s.AnnualInterestRate,
		StartDate:             s.StartDate,
		MaturityDate:          s.MaturityDate,
		LastInterestCalcDate:  s.LastInterestCalcDate,
	}
}

func (r *accountRepository) CreatePaymentAccount(ctx context.Context, userID string) (*account.Account, error) {
	accountNumber := generateAccountNumber()
	schema := &Account{
		ID:            pkg.NewUUIDV7(),
		UserID:        userID,
		AccountNumber: accountNumber,
		AccountType:   "PAYMENT",
		Balance:       0,
		Status:        "ACTIVE",
	}

	if err := r.db.WithContext(ctx).Table(AccountsTableName).Create(schema).Error; err != nil {
		return nil, err
	}

	return schema.ToDomain(), nil
}

func (r *accountRepository) CreateFixedSavingsAccount(ctx context.Context, userID string, req *account.CreateFixedSavingsAccountRequest) (*account.Account, error) {
	accountNumber := generateAccountNumber()
	schema := &Account{
		ID:            pkg.NewUUIDV7(),
		UserID:        userID,
		AccountNumber: accountNumber,
		AccountType:   "FIXED_SAVINGS",
		Balance:       0,
		Status:        "ACTIVE",
	}

	if err := r.db.WithContext(ctx).Table(AccountsTableName).Create(schema).Error; err != nil {
		return nil, err
	}

	return schema.ToDomain(), nil
}

func (r *accountRepository) CreateFlexibleSavingsAccount(ctx context.Context, userID string) (*account.Account, error) {
	accountNumber := generateAccountNumber()
	schema := &Account{
		ID:            pkg.NewUUIDV7(),
		UserID:        userID,
		AccountNumber: accountNumber,
		AccountType:   "FLEXIBLE_SAVINGS",
		Balance:       0,
		Status:        "ACTIVE",
	}

	if err := r.db.WithContext(ctx).Table(AccountsTableName).Create(schema).Error; err != nil {
		return nil, err
	}

	return schema.ToDomain(), nil
}

func (r *accountRepository) GetAccountsByUserID(ctx context.Context, userID string) ([]*account.Account, error) {
	var schemas []Account
	if err := r.db.WithContext(ctx).Table(AccountsTableName).Where("user_id = ?", userID).Find(&schemas).Error; err != nil {
		return nil, err
	}

	var accounts []*account.Account
	for _, schema := range schemas {
		accounts = append(accounts, schema.ToDomain())
	}

	return accounts, nil
}

func (r *accountRepository) GetAccountByID(ctx context.Context, accountID string) (*account.Account, error) {
	var schema Account
	if err := r.db.WithContext(ctx).Table(AccountsTableName).Where("id = ?", accountID).First(&schema).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrUserNotFound // Reuse existing error
		}
		return nil, err
	}

	return schema.ToDomain(), nil
}

func (r *accountRepository) CountPaymentAccountsByUserID(ctx context.Context, userID string) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Table(AccountsTableName).
		Where("user_id = ? AND account_type = ?", userID, "PAYMENT").
		Count(&count).Error
	return count, err
}

func (r *accountRepository) CountSavingsAccountsByUserID(ctx context.Context, userID string) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Table(AccountsTableName).
		Where("user_id = ? AND account_type IN (?, ?)", userID, "FIXED_SAVINGS", "FLEXIBLE_SAVINGS").
		Count(&count).Error
	return count, err
}

func (r *accountRepository) UpdateAccountBalance(ctx context.Context, accountID string, newBalance float64) error {
	return r.db.WithContext(ctx).Table(AccountsTableName).
		Where("id = ?", accountID).
		Update("balance", newBalance).Error
}

func (r *savingsAccountDetailRepository) CreateSavingsAccountDetail(ctx context.Context, detail *account.SavingsAccountDetail) error {
	schema := &SavingsAccountDetail{
		AccountID:             detail.AccountID,
		IsFixedTerm:           detail.IsFixedTerm,
		TermMonths:            detail.TermMonths,
		AnnualInterestRate:    detail.AnnualInterestRate,
		StartDate:             detail.StartDate,
		MaturityDate:          detail.MaturityDate,
		LastInterestCalcDate:  detail.LastInterestCalcDate,
	}

	return r.db.WithContext(ctx).Table(SavingsAccountDetailsTableName).Create(schema).Error
}

func (r *savingsAccountDetailRepository) GetSavingsAccountDetailByAccountID(ctx context.Context, accountID string) (*account.SavingsAccountDetail, error) {
	var schema SavingsAccountDetail
	if err := r.db.WithContext(ctx).Table(SavingsAccountDetailsTableName).Where("account_id = ?", accountID).First(&schema).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrUserNotFound // Reuse existing error
		}
		return nil, err
	}

	return schema.ToDomain(), nil
}

func (r *savingsAccountDetailRepository) UpdateLastInterestCalcDate(ctx context.Context, accountID string, date *time.Time) error {
	return r.db.WithContext(ctx).Table(SavingsAccountDetailsTableName).
		Where("account_id = ?", accountID).
		Update("last_interest_calc_date", date).Error
}

// generateAccountNumber generates a unique 10-digit account number
func generateAccountNumber() string {
	rand.Seed(time.Now().UnixNano())
	return fmt.Sprintf("%010d", rand.Intn(10000000000))
}