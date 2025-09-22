package user

import (
	"e-wallet/domain/user"
	"net/http"
)

type userService struct {
	client *http.Client
}

func NewUserService() user.UserService {
	return &userService{client: &http.Client{}}
}
