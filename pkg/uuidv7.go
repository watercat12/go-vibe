package pkg

import "github.com/google/uuid"

func NewUUIDV7() string {
	u,_ := uuid.NewV7()
	return u.String()
}