package account

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"time"

	"e-wallet/internal/domain/account"
	"e-wallet/internal/ports"
)

type accountService struct {
	userRepo    ports.UserRepository
	accountRepo ports.AccountRepository
	savingsRepo ports.SavingsAccountDetailRepository
}

func NewAccountService(userRepo ports.UserRepository, accountRepo ports.AccountRepository, savingsRepo ports.SavingsAccountDetailRepository) ports.AccountService {
	return &accountService{
		userRepo:    userRepo,
		accountRepo: accountRepo,
		savingsRepo: savingsRepo,
	}
}

func (s *accountService) CreatePaymentAccount(ctx context.Context, userID string) (*account.Account, error) {
	// Check if user profile is completed
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, err
	}
	if !user.IsProfileCompleted {
		return nil, errors.New("user profile must be completed before creating accounts")
	}

	// Check limit: max 1 payment account per user
	count, err := s.accountRepo.CountPaymentAccountsByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}
	if count >= 1 {
		return nil, errors.New("user can have at most 1 payment account")
	}

	// Create account
	acc, err := s.accountRepo.CreatePaymentAccount(ctx, userID)
	if err != nil {
		return nil, err
	}

	return acc, nil
}

func (s *accountService) CreateFixedSavingsAccount(ctx context.Context, userID string, req *account.CreateFixedSavingsAccountRequest) (*account.Account, error) {
	// Check if user profile is completed
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, err
	}
	if !user.IsProfileCompleted {
		return nil, errors.New("user profile must be completed before creating accounts")
	}

	// Check limit: max 5 savings accounts per user
	count, err := s.accountRepo.CountSavingsAccountsByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}
	if count >= 5 {
		return nil, errors.New("user can have at most 5 savings accounts")
	}

	// Validate term code and get interest rate
	termMonths, interestRate, err := s.getFixedSavingsTermDetails(req.TermCode)
	if err != nil {
		return nil, err
	}

	// Create account
	acc, err := s.accountRepo.CreateFixedSavingsAccount(ctx, userID, req)
	if err != nil {
		return nil, err
	}

	// Create savings detail
	maturityDate := acc.CreatedAt.AddDate(0, termMonths, 0)
	detail := &account.SavingsAccountDetail{
		AccountID:             acc.ID,
		IsFixedTerm:           true,
		TermMonths:            &termMonths,
		AnnualInterestRate:    interestRate,
		StartDate:             acc.CreatedAt,
		MaturityDate:          &maturityDate,
		LastInterestCalcDate:  nil, // Will be set on first interest calculation
	}

	if err := s.savingsRepo.CreateSavingsAccountDetail(ctx, detail); err != nil {
		return nil, err
	}

	return acc, nil
}

func (s *accountService) CreateFlexibleSavingsAccount(ctx context.Context, userID string) (*account.Account, error) {
	// Check if user profile is completed
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, err
	}
	if !user.IsProfileCompleted {
		return nil, errors.New("user profile must be completed before creating accounts")
	}

	// Check limit: max 5 savings accounts per user
	count, err := s.accountRepo.CountSavingsAccountsByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}
	if count >= 5 {
		return nil, errors.New("user can have at most 5 savings accounts")
	}

	// Create account
	acc, err := s.accountRepo.CreateFlexibleSavingsAccount(ctx, userID)
	if err != nil {
		return nil, err
	}

	// Create savings detail with promotional rate (0.8%)
	detail := &account.SavingsAccountDetail{
		AccountID:             acc.ID,
		IsFixedTerm:           false,
		TermMonths:            nil,
		AnnualInterestRate:    0.008, // 0.8% promotional rate
		StartDate:             acc.CreatedAt,
		MaturityDate:          nil,
		LastInterestCalcDate:  &acc.CreatedAt,
	}

	if err := s.savingsRepo.CreateSavingsAccountDetail(ctx, detail); err != nil {
		return nil, err
	}

	return acc, nil
}

func (s *accountService) ListAccounts(ctx context.Context, userID string) (*account.ListAccountsResponse, error) {
	accounts, err := s.accountRepo.GetAccountsByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	var response []account.AccountWithDetailsResponse
	for _, acc := range accounts {
		item := account.AccountWithDetailsResponse{
			AccountResponse: account.AccountResponse{
				ID:            acc.ID,
				UserID:        acc.UserID,
				AccountNumber: acc.AccountNumber,
				AccountType:   acc.AccountType,
				Balance:       acc.Balance,
				Status:        acc.Status,
				CreatedAt:     acc.CreatedAt,
				UpdatedAt:     acc.UpdatedAt,
			},
		}

		// Add savings details if it's a savings account
		if acc.AccountType == "FIXED_SAVINGS" || acc.AccountType == "FLEXIBLE_SAVINGS" {
			detail, err := s.savingsRepo.GetSavingsAccountDetailByAccountID(ctx, acc.ID)
			if err == nil && detail != nil {
				item.SavingsDetail = &account.SavingsAccountDetailResponse{
					AccountID:             detail.AccountID,
					IsFixedTerm:           detail.IsFixedTerm,
					TermMonths:            detail.TermMonths,
					AnnualInterestRate:    detail.AnnualInterestRate,
					StartDate:             detail.StartDate,
					MaturityDate:          detail.MaturityDate,
					LastInterestCalcDate:  detail.LastInterestCalcDate,
				}
			}
		}

		response = append(response, item)
	}

	return &account.ListAccountsResponse{Accounts: response}, nil
}

func (s *accountService) getFixedSavingsTermDetails(termCode string) (int, float64, error) {
	switch termCode {
	case "1":
		return 1, 0.006, nil // 0.6%
	case "3":
		return 3, 0.018, nil // 1.8%
	case "6":
		return 6, 0.036, nil // 3.6%
	case "8":
		return 8, 0.048, nil // 4.8%
	case "12":
		return 12, 0.072, nil // 7.2%
	default:
		return 0, 0, fmt.Errorf("invalid term code: %s", termCode)
	}
}

// generateAccountNumber generates a unique 10-digit account number
func generateAccountNumber() string {
	rand.Seed(time.Now().UnixNano())
	return fmt.Sprintf("%010d", rand.Intn(10000000000))
}