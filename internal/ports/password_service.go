package ports

type PasswordService interface {
	HashPassword(password string) (string, error)
	CheckPassword(hashedPassword, password string) error
}