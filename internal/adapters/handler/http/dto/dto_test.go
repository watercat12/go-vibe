package dto

import (
	"strings"
	"testing"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/assert"

	"e-wallet/internal/domain/account"
	"e-wallet/internal/domain/bank_link"
	"e-wallet/internal/domain/user"
)

func TestValidatePassword(t *testing.T) {
	v := validator.New()
	RegisterCustomValidations(v)

	tests := []struct {
		name     string
		password string
		want     bool
	}{
		{
			name:     "valid password",
			password: "ValidPass123!",
			want:     true,
		},
		{
			name:     "too short",
			password: "Short1!",
			want:     false,
		},
		{
			name:     "no upper",
			password: "validpass123!",
			want:     false,
		},
		{
			name:     "no lower",
			password: "VALIDPASS123!",
			want:     false,
		},
		{
			name:     "no digit",
			password: "ValidPass!",
			want:     false,
		},
		{
			name:     "no special",
			password: "ValidPass123",
			want:     false,
		},
		{
			name:     "too long",
			password: strings.Repeat("a", 51) + "1A!",
			want:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := v.Var(tt.password, "password")
			if tt.want {
				assert.NoError(t, err)
			} else {
				assert.Error(t, err)
			}
		})
	}
}

func TestValidateTeam(t *testing.T) {
	v := validator.New()
	RegisterCustomValidations(v)

	tests := []struct {
		name string
		team string
		want bool
	}{
		{
			name: "valid Front End",
			team: "Front End",
			want: true,
		},
		{
			name: "valid Back End",
			team: "Back End",
			want: true,
		},
		{
			name: "valid QA",
			team: "QA",
			want: true,
		},
		{
			name: "valid Admin",
			team: "Admin",
			want: true,
		},
		{
			name: "valid Brse",
			team: "Brse",
			want: true,
		},
		{
			name: "valid Design",
			team: "Design",
			want: true,
		},
		{
			name: "valid Others",
			team: "Others",
			want: true,
		},
		{
			name: "invalid team",
			team: "Invalid",
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := v.Var(tt.team, "team")
			if tt.want {
				assert.NoError(t, err)
			} else {
				assert.Error(t, err)
			}
		})
	}
}

func TestNewUserResponse(t *testing.T) {
	now := time.Now()
	user := &user.User{
		ID:        "user-123",
		Username:  "testuser",
		Email:     "test@example.com",
		CreatedAt: now,
		UpdatedAt: now,
	}

	resp := NewUserResponse(user)

	assert.Equal(t, user.ID, resp.ID)
	assert.Equal(t, user.Username, resp.Username)
	assert.Equal(t, user.Email, resp.Email)
	assert.Equal(t, user.CreatedAt, resp.CreatedAt)
	assert.Equal(t, user.UpdatedAt, resp.UpdatedAt)
}

func TestNewCreateUserResponse(t *testing.T) {
	resp := NewCreateUserResponse()

	assert.NotNil(t, resp)
	assert.Nil(t, resp.User)
	assert.Empty(t, resp.Token)
}

func TestNewLoginUserResponse(t *testing.T) {
	now := time.Now()
	user := &user.User{
		ID:        "user-123",
		Username:  "testuser",
		Email:     "test@example.com",
		CreatedAt: now,
		UpdatedAt: now,
	}
	token := "jwt-token"

	resp := NewLoginUserResponse(user, token)

	assert.NotNil(t, resp)
	assert.NotNil(t, resp.User)
	assert.Equal(t, user.ID, resp.User.ID)
	assert.Equal(t, user.Username, resp.User.Username)
	assert.Equal(t, user.Email, resp.User.Email)
	assert.Equal(t, user.CreatedAt, resp.User.CreatedAt)
	assert.Equal(t, user.UpdatedAt, resp.User.UpdatedAt)
	assert.Equal(t, token, resp.Token)
}

func TestNewUpdateProfileResponse(t *testing.T) {
	profile := &user.Profile{
		ID:          "profile-123",
		UserID:      "user-123",
		DisplayName: "Test User",
		AvatarURL:   "http://example.com/avatar.jpg",
		PhoneNumber: "123456789",
		NationalID:  "123456789",
		BirthYear:   1990,
		Gender:      "Male",
		Team:        "Back End",
	}

	resp := NewUpdateProfileResponse(profile)

	assert.NotNil(t, resp)
	assert.NotNil(t, resp.Profile)
	assert.Equal(t, profile.ID, resp.Profile.ID)
	assert.Equal(t, profile.UserID, resp.Profile.UserID)
	assert.Equal(t, profile.DisplayName, resp.Profile.DisplayName)
	assert.Equal(t, profile.AvatarURL, resp.Profile.AvatarURL)
	assert.Equal(t, profile.PhoneNumber, resp.Profile.PhoneNumber)
	assert.Equal(t, profile.NationalID, resp.Profile.NationalID)
	assert.Equal(t, profile.BirthYear, resp.Profile.BirthYear)
	assert.Equal(t, profile.Gender, resp.Profile.Gender)
	assert.Equal(t, profile.Team, resp.Profile.Team)
}

func TestNewProfileResponse(t *testing.T) {
	profile := &user.Profile{
		ID:          "profile-123",
		UserID:      "user-123",
		DisplayName: "Test User",
		AvatarURL:   "http://example.com/avatar.jpg",
		PhoneNumber: "123456789",
		NationalID:  "123456789",
		BirthYear:   1990,
		Gender:      "Male",
		Team:        "Back End",
	}

	resp := NewProfileResponse(profile)

	assert.Equal(t, profile.ID, resp.ID)
	assert.Equal(t, profile.UserID, resp.UserID)
	assert.Equal(t, profile.DisplayName, resp.DisplayName)
	assert.Equal(t, profile.AvatarURL, resp.AvatarURL)
	assert.Equal(t, profile.PhoneNumber, resp.PhoneNumber)
	assert.Equal(t, profile.NationalID, resp.NationalID)
	assert.Equal(t, profile.BirthYear, resp.BirthYear)
	assert.Equal(t, profile.Gender, resp.Gender)
	assert.Equal(t, profile.Team, resp.Team)
}

func TestNewAccountResponse(t *testing.T) {
	interestRate := 1.8
	fixedTermMonths := 3
	acc := &account.Account{
		ID:              "acc-123",
		UserID:          "user-123",
		AccountType:     account.PaymentAccountType,
		AccountNumber:   "PAY123456789",
		Balance:         100.0,
		InterestRate:    &interestRate,
		FixedTermMonths: &fixedTermMonths,
	}

	resp := NewAccountResponse(acc)

	assert.Equal(t, acc.ID, resp.ID)
	assert.Equal(t, acc.AccountNumber, resp.AccountNumber)
	assert.Equal(t, acc.Balance, resp.Balance)
	assert.Equal(t, acc.InterestRate, resp.InterestRate)
	assert.Equal(t, acc.FixedTermMonths, resp.FixedTermMonths)
}

func TestNewCreateAccountResponse(t *testing.T) {
	interestRate := 1.8
	fixedTermMonths := 3
	acc := &account.Account{
		ID:              "acc-123",
		UserID:          "user-123",
		AccountType:     account.PaymentAccountType,
		AccountNumber:   "PAY123456789",
		Balance:         100.0,
		InterestRate:    &interestRate,
		FixedTermMonths: &fixedTermMonths,
	}

	resp := NewCreateAccountResponse(acc)

	assert.NotNil(t, resp)
	assert.NotNil(t, resp.Account)
	assert.Equal(t, acc.ID, resp.Account.ID)
	assert.Equal(t, acc.AccountNumber, resp.Account.AccountNumber)
	assert.Equal(t, acc.Balance, resp.Account.Balance)
	assert.Equal(t, acc.InterestRate, resp.Account.InterestRate)
	assert.Equal(t, acc.FixedTermMonths, resp.Account.FixedTermMonths)
}

func TestNewListAccountsResponse(t *testing.T) {
	interestRate1 := 1.8
	fixedTermMonths1 := 3
	interestRate2 := 2.5
	fixedTermMonths2 := 6
	accounts := []*account.Account{
		{
			ID:              "acc-123",
			UserID:          "user-123",
			AccountType:     account.PaymentAccountType,
			AccountNumber:   "PAY123456789",
			AccountName:     "Payment Account",
			Balance:         100.0,
			InterestRate:    &interestRate1,
			FixedTermMonths: &fixedTermMonths1,
		},
		{
			ID:              "acc-456",
			UserID:          "user-123",
			AccountType:     account.FlexibleSavingsAccountType,
			AccountNumber:   "SAV987654321",
			AccountName:     "Savings Account",
			Balance:         500.0,
			InterestRate:    &interestRate2,
			FixedTermMonths: &fixedTermMonths2,
		},
	}

	resp := NewListAccountsResponse(accounts)

	assert.NotNil(t, resp)
	assert.Len(t, resp.Accounts, 2)
	assert.Equal(t, accounts[0].ID, resp.Accounts[0].ID)
	assert.Equal(t, accounts[0].AccountNumber, resp.Accounts[0].AccountNumber)
	assert.Equal(t, accounts[0].AccountType, resp.Accounts[0].AccountType)
	assert.Equal(t, accounts[0].AccountName, resp.Accounts[0].AccountName)
	assert.Equal(t, accounts[0].Balance, resp.Accounts[0].Balance)
	assert.Equal(t, accounts[0].InterestRate, resp.Accounts[0].InterestRate)
	assert.Equal(t, accounts[0].FixedTermMonths, resp.Accounts[0].FixedTermMonths)
	assert.Equal(t, accounts[1].ID, resp.Accounts[1].ID)
	assert.Equal(t, accounts[1].AccountNumber, resp.Accounts[1].AccountNumber)
	assert.Equal(t, accounts[1].AccountType, resp.Accounts[1].AccountType)
	assert.Equal(t, accounts[1].AccountName, resp.Accounts[1].AccountName)
	assert.Equal(t, accounts[1].Balance, resp.Accounts[1].Balance)
	assert.Equal(t, accounts[1].InterestRate, resp.Accounts[1].InterestRate)
	assert.Equal(t, accounts[1].FixedTermMonths, resp.Accounts[1].FixedTermMonths)
}

func TestNewLinkBankAccountResponse(t *testing.T) {
	bankLink := &bank_link.BankLink{
		ID:          "link-123",
		UserID:      "user-123",
		BankCode:    "BANK001",
		AccountType: "SAVINGS",
		Status:      "ACTIVE",
	}

	resp := NewLinkBankAccountResponse(bankLink)

	assert.Equal(t, bankLink.ID, resp.ID)
	assert.Equal(t, bankLink.BankCode, resp.BankCode)
	assert.Equal(t, bankLink.AccountType, resp.AccountType)
	assert.Equal(t, bankLink.Status, resp.Status)
}