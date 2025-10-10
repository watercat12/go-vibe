package postgres

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"e-wallet/internal/domain/user"
	"e-wallet/pkg"
)

func TestProfileRepository_Create(t *testing.T) {
	db := setupTestDB(t)
	userRepo := NewUserRepository(db)
	profileRepo := NewProfileRepository(db)

	tests := []struct {
		name        string
		setupUser   func() string
		profile     func(userID string) *user.Profile
		expectError bool
	}{
		{
			name: "success - create profile",
			setupUser: func() string {
				user := &user.User{
					ID:           pkg.NewUUIDV7(),
					Username:     "testuser",
					Email:        "test@example.com",
					PasswordHash: "hashedpassword",
				}
				_, err := userRepo.Create(context.Background(), user)
				require.NoError(t, err)
				return user.ID
			},
			profile: func(userID string) *user.Profile {
				return &user.Profile{
					ID:          pkg.NewUUIDV7(),
					UserID:      userID,
					DisplayName: "John Doe",
					AvatarURL:   "https://example.com/avatar.jpg",
					PhoneNumber: "1234567890",
					NationalID:  "123456789",
					BirthYear:   1990,
					Gender:      "Male",
					Team:        "Engineering",
				}
			},
			expectError: false,
		},
		{
			name: "success - create profile with minimal data",
			setupUser: func() string {
				user := &user.User{
					ID:           pkg.NewUUIDV7(),
					Username:     "testuser2",
					Email:        "test2@example.com",
					PasswordHash: "hashedpassword",
				}
				_, err := userRepo.Create(context.Background(), user)
				require.NoError(t, err)
				return user.ID
			},
			profile: func(userID string) *user.Profile {
				return &user.Profile{
					ID:          pkg.NewUUIDV7(),
					UserID:      userID,
					DisplayName: "Jane Doe",
					BirthYear:   1985,
					Gender:      "Female",
				}
			},
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			userID := tt.setupUser()
			testProfile := tt.profile(userID)
			result, err := profileRepo.Create(context.Background(), testProfile)

			if tt.expectError {
				assert.Error(t, err)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.Equal(t, testProfile.ID, result.ID)
				assert.Equal(t, testProfile.UserID, result.UserID)
				assert.Equal(t, testProfile.DisplayName, result.DisplayName)
				assert.Equal(t, testProfile.AvatarURL, result.AvatarURL)
				assert.Equal(t, testProfile.PhoneNumber, result.PhoneNumber)
				assert.Equal(t, testProfile.NationalID, result.NationalID)
				assert.Equal(t, testProfile.BirthYear, result.BirthYear)
				assert.Equal(t, testProfile.Gender, result.Gender)
				assert.Equal(t, testProfile.Team, result.Team)
				assert.NotZero(t, result.CreatedAt)
				assert.NotZero(t, result.UpdatedAt)
			}
		})
	}
}

func TestProfileRepository_Create_DBError(t *testing.T) {
	db := setupTestDB(t)
	repo := NewProfileRepository(db)

	// Close the database connection to simulate DB error
	sqlDB, _ := db.DB()
	sqlDB.Close()

	profile := &user.Profile{
		ID:          pkg.NewUUIDV7(),
		UserID:      pkg.NewUUIDV7(),
		DisplayName: "John Doe",
		BirthYear:   1990,
		Gender:      "Male",
	}

	result, err := repo.Create(context.Background(), profile)

	assert.Error(t, err)
	assert.Nil(t, result)
}

func TestProfileRepository_GetByUserID(t *testing.T) {
	db := setupTestDB(t)
	userRepo := NewUserRepository(db)
	repo := NewProfileRepository(db)

	// Setup test data
	testUser := &user.User{
		ID:           pkg.NewUUIDV7(),
		Username:     "testuser",
		Email:        "test@example.com",
		PasswordHash: "hashedpassword",
	}
	_, err := userRepo.Create(context.Background(), testUser)
	require.NoError(t, err)

	testProfile := &user.Profile{
		ID:          pkg.NewUUIDV7(),
		UserID:      testUser.ID,
		DisplayName: "John Doe",
		AvatarURL:   "https://example.com/avatar.jpg",
		PhoneNumber: "1234567890",
		NationalID:  "123456789",
		BirthYear:   1990,
		Gender:      "Male",
		Team:        "Engineering",
	}
	_, err = repo.Create(context.Background(), testProfile)
	require.NoError(t, err)

	tests := []struct {
		name        string
		userID      string
		expectError bool
		expectedErr error
	}{
		{
			name:        "success - get existing profile",
			userID:      testProfile.UserID,
			expectError: false,
		},
		{
			name:        "error - profile not found",
			userID:      pkg.NewUUIDV7(),
			expectError: true,
			expectedErr: ErrProfileNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := repo.GetByUserID(context.Background(), tt.userID)

			if tt.expectError {
				assert.Error(t, err)
				assert.Nil(t, result)
				if tt.expectedErr != nil {
					assert.Equal(t, tt.expectedErr, err)
				}
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.Equal(t, testProfile.ID, result.ID)
				assert.Equal(t, testProfile.UserID, result.UserID)
				assert.Equal(t, testProfile.DisplayName, result.DisplayName)
				assert.Equal(t, testProfile.AvatarURL, result.AvatarURL)
				assert.Equal(t, testProfile.PhoneNumber, result.PhoneNumber)
				assert.Equal(t, testProfile.NationalID, result.NationalID)
				assert.Equal(t, testProfile.BirthYear, result.BirthYear)
				assert.Equal(t, testProfile.Gender, result.Gender)
				assert.Equal(t, testProfile.Team, result.Team)
			}
		})
	}
}

func TestProfileRepository_GetByUserID_DBError(t *testing.T) {
	db := setupTestDB(t)
	repo := NewProfileRepository(db)

	// Close the database connection to simulate DB error
	sqlDB, _ := db.DB()
	sqlDB.Close()

	result, err := repo.GetByUserID(context.Background(), pkg.NewUUIDV7())

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.NotEqual(t, ErrProfileNotFound, err)
}

func TestProfileRepository_Update(t *testing.T) {
	db := setupTestDB(t)
	userRepo := NewUserRepository(db)
	repo := NewProfileRepository(db)

	// Setup test data
	testUser := &user.User{
		ID:           pkg.NewUUIDV7(),
		Username:     "testuser",
		Email:        "test@example.com",
		PasswordHash: "hashedpassword",
	}
	_, err := userRepo.Create(context.Background(), testUser)
	require.NoError(t, err)

	testProfile := &user.Profile{
		ID:          pkg.NewUUIDV7(),
		UserID:      testUser.ID,
		DisplayName: "John Doe",
		AvatarURL:   "https://example.com/avatar.jpg",
		PhoneNumber: "1234567890",
		NationalID:  "123456789",
		BirthYear:   1990,
		Gender:      "Male",
		Team:        "Engineering",
	}
	_, err = repo.Create(context.Background(), testProfile)
	require.NoError(t, err)

	// Update profile
	updatedProfile := &user.Profile{
		UserID:      testUser.ID,
		DisplayName: "John Smith",
		AvatarURL:   "https://example.com/new-avatar.jpg",
		PhoneNumber: "0987654321",
		NationalID:  "987654321",
		BirthYear:   1985,
		Gender:      "Male",
		Team:        "Management",
	}

	result, err := repo.Update(context.Background(), updatedProfile)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, testProfile.ID, result.ID) // ID should remain the same
	assert.Equal(t, updatedProfile.UserID, result.UserID)
	assert.Equal(t, updatedProfile.DisplayName, result.DisplayName)
	assert.Equal(t, updatedProfile.AvatarURL, result.AvatarURL)
	assert.Equal(t, updatedProfile.PhoneNumber, result.PhoneNumber)
	assert.Equal(t, updatedProfile.NationalID, result.NationalID)
	assert.Equal(t, updatedProfile.BirthYear, result.BirthYear)
	assert.Equal(t, updatedProfile.Gender, result.Gender)
	assert.Equal(t, updatedProfile.Team, result.Team)
	assert.NotZero(t, result.UpdatedAt)
}

func TestProfileRepository_Update_DBError(t *testing.T) {
	db := setupTestDB(t)
	repo := NewProfileRepository(db)

	// Close the database connection to simulate DB error
	sqlDB, _ := db.DB()
	sqlDB.Close()

	profile := &user.Profile{
		UserID:      pkg.NewUUIDV7(),
		DisplayName: "John Doe",
		BirthYear:   1990,
		Gender:      "Male",
	}

	result, err := repo.Update(context.Background(), profile)

	assert.Error(t, err)
	assert.Nil(t, result)
}

func TestProfileRepository_Upsert(t *testing.T) {
	db := setupTestDB(t)
	userRepo := NewUserRepository(db)
	repo := NewProfileRepository(db)

	// Setup test data
	testUser := &user.User{
		ID:           pkg.NewUUIDV7(),
		Username:     "testuser",
		Email:        "test@example.com",
		PasswordHash: "hashedpassword",
	}
	_, err := userRepo.Create(context.Background(), testUser)
	require.NoError(t, err)

	t.Run("insert new profile", func(t *testing.T) {
		profile := &user.Profile{
			ID:          pkg.NewUUIDV7(),
			UserID:      testUser.ID,
			DisplayName: "John Doe",
			AvatarURL:   "https://example.com/avatar.jpg",
			PhoneNumber: "1234567890",
			NationalID:  "123456789",
			BirthYear:   1990,
			Gender:      "Male",
			Team:        "Engineering",
		}

		result, err := repo.Upsert(context.Background(), profile)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, profile.ID, result.ID)
		assert.Equal(t, profile.UserID, result.UserID)
		assert.Equal(t, profile.DisplayName, result.DisplayName)
		assert.Equal(t, profile.AvatarURL, result.AvatarURL)
		assert.Equal(t, profile.PhoneNumber, result.PhoneNumber)
		assert.Equal(t, profile.NationalID, result.NationalID)
		assert.Equal(t, profile.BirthYear, result.BirthYear)
		assert.Equal(t, profile.Gender, result.Gender)
		assert.Equal(t, profile.Team, result.Team)
	})

	t.Run("update existing profile", func(t *testing.T) {
		updatedProfile := &user.Profile{
			ID:          pkg.NewUUIDV7(),
			UserID:      testUser.ID,
			DisplayName: "John Smith",
			AvatarURL:   "https://example.com/new-avatar.jpg",
			PhoneNumber: "0987654321",
			NationalID:  "987654321",
			BirthYear:   1985,
			Gender:      "Male",
			Team:        "Management",
		}

		result, err := repo.Upsert(context.Background(), updatedProfile)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, updatedProfile.ID, result.ID)
		assert.Equal(t, updatedProfile.UserID, result.UserID)
		assert.Equal(t, updatedProfile.DisplayName, result.DisplayName)
		assert.Equal(t, updatedProfile.AvatarURL, result.AvatarURL)
		assert.Equal(t, updatedProfile.PhoneNumber, result.PhoneNumber)
		assert.Equal(t, updatedProfile.NationalID, result.NationalID)
		assert.Equal(t, updatedProfile.BirthYear, result.BirthYear)
		assert.Equal(t, updatedProfile.Gender, result.Gender)
		assert.Equal(t, updatedProfile.Team, result.Team)
	})
}

func TestProfileRepository_Upsert_DBError(t *testing.T) {
	db := setupTestDB(t)
	repo := NewProfileRepository(db)

	// Close the database connection to simulate DB error
	sqlDB, _ := db.DB()
	sqlDB.Close()

	profile := &user.Profile{
		ID:          pkg.NewUUIDV7(),
		UserID:      pkg.NewUUIDV7(),
		DisplayName: "John Doe",
		BirthYear:   1990,
		Gender:      "Male",
	}

	result, err := repo.Upsert(context.Background(), profile)

	assert.Error(t, err)
	assert.Nil(t, result)
}