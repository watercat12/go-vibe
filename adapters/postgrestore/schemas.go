package postgrestore

import (
	"e-wallet/domain/user"
	"time"
)

const (
	UsersTableName     = "users"
)

type User struct {
	ID        int       `gorm:"primaryKey"`
	Email     string    `gorm:"uniqueIndex;not null"`
	Name      string    `gorm:"not null"`
	Password  string    `gorm:"not null"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
}

func (User) TableName() string {
	return "users"
}

// ToDomain converts schema to domain model
func (u *User) ToDomain() *user.User {
	return &user.User{
		ID:        u.ID,
		Email:     u.Email,
		Name:      u.Name,
		Password:  u.Password,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	}
}