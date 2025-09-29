package dto

type UpdateProfileRequest struct {
	Username    string `json:"username" validate:"required"`
	DisplayName string `json:"display_name"`
	AvatarURL   string `json:"avatar_url"`
	PhoneNumber string `json:"phone_number"`
	NationalID  string `json:"national_id"`
	BirthYear   int    `json:"birth_year"`
	Gender      string `json:"gender"`
	Team        string `json:"team" validate:"team"`
}