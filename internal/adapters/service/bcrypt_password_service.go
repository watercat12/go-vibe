package service

import (
	"golang.org/x/crypto/bcrypt"
)

type bcryptPasswordService struct{}

func NewBcryptPasswordService() *bcryptPasswordService {
	return &bcryptPasswordService{}
}

func (s *bcryptPasswordService) HashPassword(password string) (string, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashed), nil
}

func (s *bcryptPasswordService) CheckPassword(hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}