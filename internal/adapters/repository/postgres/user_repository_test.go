package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"

	"e-wallet/internal/config"
	"e-wallet/internal/domain/user"
	"e-wallet/pkg"
	_ "github.com/lib/pq"
	migrate "github.com/rubenv/sql-migrate"
)

func TestMain(m *testing.M) {
	// Set up test database environment variables
	os.Setenv("DB_HOST", "pi.local")
	os.Setenv("DB_PORT", "54321")
	os.Setenv("DB_USER", "postgres")
	os.Setenv("DB_PASS", "123456")
	os.Setenv("DB_NAME", "e_wallet_test")
	os.Setenv("ENABLE_SSL", "false")

	code := m.Run()

	os.Exit(code)
}

func setupTestDB(t *testing.T) *gorm.DB {
	cfg, err := config.LoadConfig()
	require.NoError(t, err)

	opts := ParseFromConfig(cfg)

	// Generate unique database name
	uniqueDBName := fmt.Sprintf("test_db_%d", time.Now().UnixNano())

	// Create admin DSN to connect to postgres database
	sslmode := "disable"
	if opts.SSLMode {
		sslmode = "enable"
	}
	adminDSN := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=postgres sslmode=%s",
		opts.Host, opts.Port, opts.DBUser, opts.Password, sslmode,
	)

	// Connect to postgres database to create new database
	adminDB, err := sql.Open("postgres", adminDSN)
	require.NoError(t, err)
	defer adminDB.Close()

	// Create the unique database
	_, err = adminDB.Exec("CREATE DATABASE " + uniqueDBName)
	require.NoError(t, err)

	// Set the database name to the unique one
	opts.DBName = uniqueDBName

	// Connect to the new database
	db, err := NewConnection(opts)
	require.NoError(t, err)

	dbTest,err :=db.DB()
	require.NoError(t, err)
	migrations := &migrate.FileMigrationSource{
		Dir: "../../../../migrations",
	}
	_, err = migrate.Exec(dbTest, "postgres", migrations, migrate.Up)
	require.NoError(t, err)

	// Cleanup: drop the database after test
	t.Cleanup(func() {
		// Close the GORM connection
		sqlDB, _ := db.DB()
		sqlDB.Close()

		// Connect again to postgres to drop the database
		adminDB2, err := sql.Open("postgres", adminDSN)
		if err == nil {
			defer adminDB2.Close()
			adminDB2.Exec("DROP DATABASE IF EXISTS " + uniqueDBName)
		}
	})

	return db
}

func TestUserRepository_Create(t *testing.T) {
	db := setupTestDB(t)
	repo := NewUserRepository(db)

	tests := []struct {
		name        string
		user        *user.User
		expectError bool
	}{
		{
			name: "success - create user",
			user: &user.User{
				ID:           pkg.NewUUIDV7(),
				Username:     "testuser",
				Email:        "test@example.com",
				PasswordHash: "hashedpassword",
			},
			expectError: false,
		},
		{
			name: "success - create user with email verified",
			user: &user.User{
				ID:              pkg.NewUUIDV7(),
				Username:        "testuser2",
				Email:           "test2@example.com",
				PasswordHash:    "hashedpassword",
				IsEmailVerified: true,
			},
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := repo.Create(context.Background(), tt.user)

			if tt.expectError {
				assert.Error(t, err)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.Equal(t, tt.user.ID, result.ID)
				assert.Equal(t, tt.user.Username, result.Username)
				assert.Equal(t, tt.user.Email, result.Email)
				assert.Equal(t, tt.user.PasswordHash, result.PasswordHash)
				assert.Equal(t, tt.user.IsEmailVerified, result.IsEmailVerified)
				assert.NotZero(t, result.CreatedAt)
				assert.NotZero(t, result.UpdatedAt)
			}
		})
	}
}

func TestUserRepository_Create_DBError(t *testing.T) {
	db := setupTestDB(t)
	repo := NewUserRepository(db)

	// Close the database connection to simulate DB error
	sqlDB, _ := db.DB()
	sqlDB.Close()

	user := &user.User{
		ID:           pkg.NewUUIDV7(),
		Username:     "testuser",
		Email:        "test@example.com",
		PasswordHash: "hashedpassword",
	}

	result, err := repo.Create(context.Background(), user)

	assert.Error(t, err)
	assert.Nil(t, result)
}

func TestUserRepository_GetByEmail_DBError(t *testing.T) {
	db := setupTestDB(t)
	repo := NewUserRepository(db)

	// Close the database connection to simulate DB error
	sqlDB, _ := db.DB()
	sqlDB.Close()

	result, err := repo.GetByEmail(context.Background(), "test@example.com")

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.NotEqual(t, ErrUserNotFound, err)
}

func TestUserRepository_GetByID_DBError(t *testing.T) {
	db := setupTestDB(t)
	repo := NewUserRepository(db)

	// Close the database connection to simulate DB error
	sqlDB, _ := db.DB()
	sqlDB.Close()

	result, err := repo.GetByID(context.Background(), pkg.NewUUIDV7())

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.NotEqual(t, ErrUserNotFound, err)
}

func TestUserRepository_GetByEmail(t *testing.T) {
	db := setupTestDB(t)
	repo := NewUserRepository(db)

	// Setup test data
	testUser := &user.User{
		ID:           pkg.NewUUIDV7(),
		Username:     "testuser",
		Email:        "test@example.com",
		PasswordHash: "hashedpassword",
	}
	_, err := repo.Create(context.Background(), testUser)
	require.NoError(t, err)

	tests := []struct {
		name        string
		email       string
		expectError bool
		expectedErr error
	}{
		{
			name:        "success - get existing user",
			email:       "test@example.com",
			expectError: false,
		},
		{
			name:        "error - user not found",
			email:       "nonexistent@example.com",
			expectError: true,
			expectedErr: ErrUserNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := repo.GetByEmail(context.Background(), tt.email)

			if tt.expectError {
				assert.Error(t, err)
				assert.Nil(t, result)
				if tt.expectedErr != nil {
					assert.Equal(t, tt.expectedErr, err)
				}
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.Equal(t, testUser.ID, result.ID)
				assert.Equal(t, testUser.Username, result.Username)
				assert.Equal(t, testUser.Email, result.Email)
				assert.Equal(t, testUser.PasswordHash, result.PasswordHash)
			}
		})
	}
}

func TestUserRepository_GetByID(t *testing.T) {
	db := setupTestDB(t)
	repo := NewUserRepository(db)

	// Setup test data
	testUser := &user.User{
		ID:           pkg.NewUUIDV7(),
		Username:     "testuser",
		Email:        "test@example.com",
		PasswordHash: "hashedpassword",
	}
	_, err := repo.Create(context.Background(), testUser)
	require.NoError(t, err)

	tests := []struct {
		name        string
		id          string
		expectError bool
		expectedErr error
	}{
		{
			name:        "success - get existing user",
			id:          testUser.ID,
			expectError: false,
		},
		{
			name:        "error - user not found",
			id:          pkg.NewUUIDV7(),
			expectError: true,
			expectedErr: ErrUserNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := repo.GetByID(context.Background(), tt.id)

			if tt.expectError {
				assert.Error(t, err)
				assert.Nil(t, result)
				if tt.expectedErr != nil {
					assert.Equal(t, tt.expectedErr, err)
				}
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.Equal(t, testUser.ID, result.ID)
				assert.Equal(t, testUser.Username, result.Username)
				assert.Equal(t, testUser.Email, result.Email)
				assert.Equal(t, testUser.PasswordHash, result.PasswordHash)
			}
		})
	}
}