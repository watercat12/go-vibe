package dto

type CreateFixedSavingsAccountRequest struct {
	TermCode string `json:"term_code" validate:"required,oneof=1 3 6 8 12" example:"12"`
}