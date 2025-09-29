package account

import (
	"context"
	"fmt"

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
	accountNumber := fmt.Sprintf("PAY%s", pkg.NewUUIDV7()[:8])

	// Create account
	acc := &account.Account{
		ID:            pkg.NewUUIDV7(),
		UserID:        userID,
		AccountType:   "payment",
		AccountNumber: accountNumber,
		Balance:       0.0,
	}

	return s.repo.Create(ctx, acc)
}