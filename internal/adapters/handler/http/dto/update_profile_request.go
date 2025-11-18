package dto

type UpdateProfileRequest struct {
	DisplayName string `json:"display_name" validate:"required"`
	AvatarURL   *string `json:"avatar_url"`
	PhoneNumber string `json:"phone_number" validate:"required,phone"`
	NationalID  string `json:"national_id" validate:"required"`
	BirthYear   int    `json:"birth_year" validate:"required,min=1900,max=2100"`
	Gender      string `json:"gender" validate:"required,oneof=MALE FEMALE OTHER"`
	Team        string `json:"team" validate:"required"`
}