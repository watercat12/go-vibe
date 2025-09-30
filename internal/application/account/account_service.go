package account

import (
	"context"
	"fmt"
	"time"

	"e-wallet/internal/domain/account"
	"e-wallet/internal/ports"
	"e-wallet/pkg"
)

type accountService struct {
	repo        ports.AccountRepository
	userRepo    ports.UserRepository
	profileRepo ports.ProfileRepository
}

func NewAccountService(repo ports.AccountRepository, userRepo ports.UserRepository, profileRepo ports.ProfileRepository) account.AccountService {
	return &accountService{repo: repo, userRepo: userRepo, profileRepo: profileRepo}
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