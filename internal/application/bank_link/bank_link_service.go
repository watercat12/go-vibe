package bank_link

import (
	"context"
	"errors"
	"e-wallet/internal/domain/bank_link"
	"e-wallet/internal/ports"
)

type bankLinkService struct {
	repo     ports.BankLinkRepository
	client   ports.BankLinkClient
}

func NewBankLinkService(repo ports.BankLinkRepository, client ports.BankLinkClient) ports.BankLinkService {
	return &bankLinkService{repo: repo, client: client}
}

func (s *bankLinkService) LinkBankAccount(ctx context.Context, userID, bankCode, accountType string) (*bank_link.BankLink, error) {
	// Check limit: max 5 bank links per user
	count, err := s.repo.CountByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}
	if count >= 5 {
		return nil, errors.New("you have linked the maximum number of 5 bank accounts")
	}

	// Call external API
	linkedAccount, err := s.client.LinkAccount(ctx, bankCode, accountType)
	if err != nil {
		return nil, err
	}

	// Create bank link with userID
	bankLink := bank_link.NewBankLink(userID, bankCode, accountType, linkedAccount.AccessToken, linkedAccount.RefreshToken, linkedAccount.ExpiresIn)

	// Save to database
	return s.repo.Create(ctx, bankLink)
}