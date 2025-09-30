package account

import "errors"

var (
	ErrLimitPaymentAccount  = errors.New("limit payment account")
	ErrLimitSavingsAccount  = errors.New("limit savings account")
	ErrInvalidTermMonths    = errors.New("invalid term months")
)