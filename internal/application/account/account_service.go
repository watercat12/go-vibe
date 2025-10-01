package account

import (
	"context"
	"fmt"
	"time"

	"e-wallet/internal/domain/account"
	"e-wallet/internal/domain/interest_history"
	"e-wallet/internal/domain/transaction"
	"e-wallet/internal/ports"
	"e-wallet/pkg"
)

type accountService struct {
	repo        ports.AccountRepository
	userRepo    ports.UserRepository
	profileRepo ports.ProfileRepository
	txRepo      ports.TransactionRepository
	ihRepo      ports.InterestHistoryRepository
}

func NewAccountService(repo ports.AccountRepository, userRepo ports.UserRepository, profileRepo ports.ProfileRepository, txRepo ports.TransactionRepository, ihRepo ports.InterestHistoryRepository) ports.AccountService {
	return &accountService{repo: repo, userRepo: userRepo, profileRepo: profileRepo, txRepo: txRepo, ihRepo: ihRepo}
}

func (s *accountService) CreatePaymentAccount(ctx context.Context, userID string) (*account.Account, error) {
	// Check if user exists
	_, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	// Check profile status (assume verified if profile exists)
	_, err = s.profileRepo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	// Check if user already has a payment account
	_, err = s.repo.GetByUserID(ctx, userID)
	if err == nil {
		return nil, ErrLimitPaymentAccount
	}

	// Generate account number (simple: PAY + uuid prefix)
	accountNumber := fmt.Sprintf("PAY%d", time.Now().UnixNano())

	// Create account
	acc := &account.Account{
		ID:            pkg.NewUUIDV7(),
		UserID:        userID,
		AccountType:   account.PaymentAccountType,
		AccountNumber: accountNumber,
		Balance:       0.0,
	}

	return s.repo.Create(ctx, acc)
}

func (s *accountService) CalculateDailyInterest(ctx context.Context) error {
	// Get yesterday's date
	yesterday := time.Now().AddDate(0, 0, -1)

	// Get all flexible savings accounts
	accounts, err := s.repo.GetFlexibleSavingsAccounts(ctx)
	if err != nil {
		return err
	}

	for _, acc := range accounts {
		// Calculate interest
		interestAmount := s.calculateInterest(acc.Balance, acc.CreatedAt)

		if interestAmount == 0 {
			continue
		}

		// Update balance
		newBalance := acc.Balance + interestAmount
		if err := s.repo.UpdateBalance(ctx, acc.ID, newBalance); err != nil {
			return err
		}

		// Create transaction
		tx := &transaction.Transaction{
			ID:              pkg.NewUUIDV7(),
			AccountID:       acc.ID,
			TransactionType: transaction.TransactionTypeInterest,
			Amount:          interestAmount,
			Status:          transaction.TransactionStatusSuccess,
			BalanceAfter:    newBalance,
		}
		if _, err := s.txRepo.Create(ctx, tx); err != nil {
			return err
		}

		// Create interest history
		ih := &interest_history.InterestHistory{
			ID:             pkg.NewUUIDV7(),
			AccountID:      acc.ID,
			Date:           yesterday,
			InterestAmount: interestAmount,
		}
		if _, err := s.ihRepo.Create(ctx, ih); err != nil {
			return err
		}
	}

	return nil
}

func (s *accountService) calculateInterest(balance float64, date time.Time) float64 {
	var annualRate float64
	switch {
	case time.Since(date) < 30*24*time.Hour:
		annualRate = 0.008
	case balance < 10000000: // under 10M
		annualRate = 0.003
	case balance < 50000000: // 10-50M
		annualRate = 0.004
	default: // over 50M
		annualRate = 0.005
	}

	dailyRate := annualRate / 365
	return balance * dailyRate
}

func (s *accountService) CreateFlexibleSavingsAccount(ctx context.Context, userID string) (*account.Account, error) {
	// Check if user exists
	_, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	// Check profile status (assume verified if profile exists)
	_, err = s.profileRepo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	// Check limit savings accounts
	count, err := s.repo.CountSavingsAccounts(ctx, userID)
	if err != nil {
		return nil, err
	}
	if count >= 5 {
		return nil, ErrLimitSavingsAccount
	}

	// Generate account number
	accountNumber := fmt.Sprintf("SAV%d", time.Now().UnixNano())

	// Create account
	acc := &account.Account{
		ID:            pkg.NewUUIDV7(),
		UserID:        userID,
		AccountType:   account.FlexibleSavingsAccountType,
		AccountNumber: accountNumber,
		Balance:       0.0,
	}

	return s.repo.Create(ctx, acc)
}

func (s *accountService) CreateFixedSavingsAccount(ctx context.Context, userID string, termMonths int) (*account.Account, error) {
	// Check if user exists
	_, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	// Check profile status (assume verified if profile exists)
	_, err = s.profileRepo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	// Validate term months
	validTerms := map[int]float64{
		1:  0.6,
		3:  1.8,
		6:  3.6,
		8:  4.8,
		12: 7.2,
	}
	interestRate, ok := validTerms[termMonths]
	if !ok {
		return nil, ErrInvalidTermMonths
	}

	// Check limit savings accounts
	count, err := s.repo.CountSavingsAccounts(ctx, userID)
	if err != nil {
		return nil, err
	}
	if count >= 5 {
		return nil, ErrLimitSavingsAccount
	}

	// Generate account number
	accountNumber := fmt.Sprintf("SAV%d", time.Now().UnixNano())

	// Create account
	acc := &account.Account{
		ID:              pkg.NewUUIDV7(),
		UserID:          userID,
		AccountType:     account.FixedSavingsAccountType,
		AccountNumber:   accountNumber,
		Balance:         0.0,
		InterestRate:    &interestRate,
		FixedTermMonths: &termMonths,
	}

	return s.repo.Create(ctx, acc)
}