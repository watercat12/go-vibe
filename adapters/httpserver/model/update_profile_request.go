package model

type UpdateProfileRequest struct {
	Name      string `json:"name"`
	Email     string `json:"email"`
	Avatar    string `json:"avatar"`
	Phone     string `json:"phone"`
	IDNumber  string `json:"id_number"`
	BirthYear int    `json:"birth_year" validate:"min=1900,max=2100"`
	Gender    string `json:"gender" validate:"oneof=male female other"`
	Team      string `json:"team" validate:"oneof=Front End Back End QA Admin Brse Design Others"`
}