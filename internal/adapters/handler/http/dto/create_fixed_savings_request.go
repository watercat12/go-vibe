package dto

type CreateFixedSavingsRequest struct {
	TermMonths int `json:"term_months" validate:"required,oneof=1 3 6 8 12"`
}