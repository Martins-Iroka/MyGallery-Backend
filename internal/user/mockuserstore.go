package user

import (
	"context"
	"time"

	"github.com/Martins-Iroka/MyGallery-Backend/internal/auth"
)

type MockUserStore struct{}

func (s *MockUserStore) ActivateUser(ctx context.Context, token string) error {
	return nil
}

func (s *MockUserStore) CreateUserAndVerificationToken(ctx context.Context, user *User, token string) error {
	return nil
}

func (s *MockUserStore) DeleteUser(ctx context.Context, userID int64) error {
	return nil
}

func (s *MockUserStore) GetUserByEmail(ctx context.Context, email string) (*User, error) {
	hashedPassword, _ := auth.HashPassword("12345")
	return &User{ID: 1, Password: hashedPassword}, nil
}

func (s *MockUserStore) GetUserByID(ctx context.Context, userID int64) (*User, error) {
	return nil, nil
}

func (s *MockUserStore) CreateRefreshToken(ctx context.Context, userID int64, tokenHash string, expiresAt time.Time) error {
	return nil
}

func (s *MockUserStore) GetUserByRefreshToken(ctx context.Context, tokenHash string) (*User, error) {
	return nil, nil
}

func (s *MockUserStore) RevokeRefreshToken(ctx context.Context, tokenHash string) error {
	return nil
}

func (s *MockUserStore) DeleteExpiredRefreshTokens(context.Context) error {
	return nil
}
