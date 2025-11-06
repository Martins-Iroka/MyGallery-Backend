package user

import (
	"context"
	"database/sql"

	"github.com/Martins-Iroka/MyGallery-Backend/internal/util"
)

type User struct {
	ID         int64  `json:"id"`
	Username   string `json:"username"`
	Email      string `json:"email"`
	Password   string `json:"-"`
	IsVerified string `json:"is_verified"`
	CreatedAt  string `json:"created_at"`
}

type UserStore struct {
	Db *sql.DB
}

func (s *UserStore) ActivateUser(ctx context.Context, token string) error {
	return nil
}

func (s *UserStore) CreateUser(ctx context.Context, tx *sql.Tx, user *User) error {
	query := `
		INSERT INTO users (email, username, password)
		VALUES ($1, $2, $3)
	`
	ctx, cancel := context.WithTimeout(ctx, util.QueryTimeoutDuration)
	defer cancel()

	_, err := tx.ExecContext(ctx, query, user.Email, user.Username, user.Password)
	if err != nil {
		switch {
		case err.Error() == `pq: duplicate key value violates unique constraint "users_email_key"`:
			return util.ErrorDuplicateEmail
		case err.Error() == `pq: duplicate key value violates unique constraint "users_username_key"`:
			return util.ErrorDuplicateUsername
		default:
			return err
		}
	}
	return nil
}

func (s *UserStore) CreateAndInviteUser(ctx context.Context, user *User, token string) error {
	return util.WithTransaction(s.Db, ctx, func(tx *sql.Tx) error {
		if err := s.CreateUser(ctx, tx, user); err != nil {
			return err
		}

		if err := s.createUserInvitation(ctx, tx, token, user.ID); err != nil {
			return err
		}

		return nil

	})
}

func (s *UserStore) DeleteUser(context.Context, int64) error {
	return nil
}

func (s *UserStore) GetUserByEmail(ctx context.Context, email string) (*User, error) {
	return nil, nil
}

func (s *UserStore) createUserInvitation(ctx context.Context, tx *sql.Tx, token string, userID int64) error {
	query := `INSERT INTO users_verification_tracking (token, user_id) VALUES ($1, $2)`

	ctx, cancel := context.WithTimeout(ctx, util.QueryTimeoutDuration)
	defer cancel()

	_, err := tx.ExecContext(ctx, query, token, userID)
	if err != nil {
		return err
	}

	return nil
}
