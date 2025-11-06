package user

import (
	"context"
	"database/sql"
)

type User struct {
	ID        int64  `json:"id"`
	Username  string `json:"username"`
	Email     string `json:"email"`
	Password  string `json:"-"`
	IsActive  string `json:"is_verified"`
	CreatedAt string `json:"created_at"`
}

type UserStore struct {
	Db *sql.DB
}

func (s *UserStore) ActivateUser(ctx context.Context, token string) error {
	return nil
}

func (s *UserStore) CreateUser(ctx context.Context, tx *sql.Tx, user *User) error {
	return nil
}

func (s *UserStore) CreateAndInviteUser(ctx context.Context, user *User, token string) error {
	return nil
}

func (s *UserStore) DeleteUser(context.Context, int64) error {
	return nil
}

func (s *UserStore) GetUserByEmail(ctx context.Context, email string) (*User, error) {
	return nil, nil
}
