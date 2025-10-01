package ports

import (
	"context"
	"e-wallet/internal/domain/interest_history"
)

type InterestHistoryRepository interface {
	Create(ctx context.Context, ih *interest_history.InterestHistory) (*interest_history.InterestHistory, error)
}