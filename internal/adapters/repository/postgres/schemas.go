package postgres

import (
	"e-wallet/internal/domain/account"
	"e-wallet/internal/domain/user"
	"time"
)

const (
	UsersTableName    = "users"
	ProfilesTableName = "profiles"
	AccountsTableName = "accounts"
)

type User struct {
	ID              string
	Username        string
	Email           string
	PasswordHash    string
	IsEmailVerified bool
	CreatedAt       time.Time `gorm:"autoCreateTime"`
	UpdatedAt       time.Time `gorm:"autoUpdateTime"`
}

func (User) TableName() string {
	return "users"
}

func (u *User) ToDomain() *user.User {
	return &user.User{
		ID:              u.ID,
		Username:        u.Username,
		Email:           u.Email,
		PasswordHash:    u.PasswordHash,
		IsEmailVerified: u.IsEmailVerified,
		CreatedAt:       u.CreatedAt,
		UpdatedAt:       u.UpdatedAt,
	}
}

type Profile struct {
	ID          string
	UserID      string
	DisplayName string
	AvatarURL   string
	PhoneNumber string
	NationalID  string
	BirthYear   int
	Gender      string
	Team        string
	CreatedAt   time.Time `gorm:"autoCreateTime"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime"`
}

func (p *Profile) ToDomain() *user.Profile {
	return &user.Profile{
		ID:          p.ID,
		UserID:      p.UserID,
		DisplayName: p.DisplayName,
		AvatarURL:   p.AvatarURL,
		PhoneNumber: p.PhoneNumber,
		NationalID:  p.NationalID,
		BirthYear:   p.BirthYear,
		Gender:      p.Gender,
		Team:        p.Team,
		CreatedAt:   p.CreatedAt,
		UpdatedAt:   p.UpdatedAt,
	}
}

type Account struct {
	ID              string
	UserID          string
	AccountType     string
	AccountNumber   string
	AccountName     string
	Balance         float64
	InterestRate    *float64
	FixedTermMonths *int
	CreatedAt       time.Time `gorm:"autoCreateTime"`
	UpdatedAt       time.Time `gorm:"autoUpdateTime"`
}

func (a *Account) ToDomain() *account.Account {
	return &account.Account{
		ID:              a.ID,
		UserID:          a.UserID,
		AccountType:     a.AccountType,
		AccountNumber:   a.AccountNumber,
		AccountName:     a.AccountName,
		Balance:         a.Balance,
		InterestRate:    a.InterestRate,
		FixedTermMonths: a.FixedTermMonths,
		CreatedAt:       a.CreatedAt,
		UpdatedAt:       a.UpdatedAt,
	}
}
