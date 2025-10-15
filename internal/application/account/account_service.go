package account

import (
	"context"
	"time"

	"e-wallet/internal/domain/account"
	"e-wallet/internal/domain/interest_history"
	"e-wallet/internal/domain/transaction"
	"e-wallet/internal/ports"
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
	_, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	_, err = s.profileRepo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	_, err = s.repo.GetPaymentAccount(ctx, userID)
	if err == nil {
		return nil, ErrLimitPaymentAccount
	}

	acc := account.NewPaymentAccount(userID)

	return s.repo.Create(ctx, acc)
}

func (s *accountService) CalculateDailyInterest(ctx context.Context) error {
	yesterday := time.Now().AddDate(0, 0, -1)

	accounts, err := s.repo.GetFlexibleSavingsAccounts(ctx)
	if err != nil {
		return err
	}

	for _, acc := range accounts {
		interestAmount := acc.CalculateDailyInterest()

		if interestAmount == 0 {
			continue
		}

		newBalance := acc.Balance + interestAmount
		if err := s.repo.UpdateBalance(ctx, acc.ID, newBalance); err != nil {
			return err
		}

		tx := transaction.NewInterestTransaction(acc.ID, interestAmount, newBalance)
		if _, err := s.txRepo.Create(ctx, tx); err != nil {
			return err
		}

		ih := interest_history.NewInterestHistory(acc.ID, yesterday, interestAmount)
		if _, err := s.ihRepo.Create(ctx, ih); err != nil {
			return err
		}
	}

	return nil
}

func (s *accountService) CreateFlexibleSavingsAccount(ctx context.Context, userID string) (*account.Account, error) {
	_, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	_, err = s.profileRepo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	count, err := s.repo.CountSavingsAccounts(ctx, userID)
	if err != nil {
		return nil, err
	}
	accountPolicy := account.NewAccountPolicy()
	if accountPolicy.LimitSavingAccount(int(count)) {
		return nil, ErrLimitSavingsAccount
	}

	acc := account.NewFlexibleSavingsAccount(userID)

	return s.repo.Create(ctx, acc)
}

func (s *accountService) CreateFixedSavingsAccount(ctx context.Context, userID string, termMonths int) (*account.Account, error) {
	_, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	_, err = s.profileRepo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	interestRate, err := account.ValidateTermMonths(termMonths)
	if err != nil {
		return nil, ErrInvalidTermMonths
	}

	count, err := s.repo.CountSavingsAccounts(ctx, userID)
	if err != nil {
		return nil, err
	}
	accountPolicy := account.NewAccountPolicy()
	if accountPolicy.LimitSavingAccount(int(count)) {
		return nil, ErrLimitSavingsAccount
	}

	acc := account.NewFixedSavingsAccount(userID, termMonths, interestRate)

	return s.repo.Create(ctx, acc)
}

func (s *accountService) GetAccountsByUserID(ctx context.Context, userID string) ([]*account.Account, error) {
	_, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	return s.repo.GetAccountsByUserID(ctx, userID)
}