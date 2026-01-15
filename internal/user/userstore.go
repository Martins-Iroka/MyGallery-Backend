package user

import (
	"context"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"log"
	"time"

	"github.com/Martins-Iroka/MyGallery-Backend/internal/util"
)

type User struct {
	ID         int64  `json:"id"`
	Username   string `json:"username"`
	Email      string `json:"email"`
	Password   string `json:"-"`
	IsVerified bool   `json:"is_verified"`
	CreatedAt  string `json:"created_at"`
}

type UserStore struct {
	Db *sql.DB
}

func (s *UserStore) ActivateUser(ctx context.Context, token string) error {
	return util.WithTransaction(s.Db, ctx, func(tx *sql.Tx) error {
		user, err := s.getUserByVerificationToken(ctx, tx, token)
		if err != nil {
			return err
		}

		if err := s.updateUser(ctx, tx, user); err != nil {
			return err
		}

		if err := s.deleteUserVerificationToken(ctx, tx, user.ID); err != nil {
			return err
		}

		return nil
	})
}

func (s *UserStore) CreateUserAndVerificationToken(ctx context.Context, user *User, token string) error {
	return util.WithTransaction(s.Db, ctx, func(tx *sql.Tx) error {
		userID, err := s.getUnverifiedUser(ctx, tx, user.Email)
		if err == nil {
			if err := s.deleteUser(ctx, tx, userID); err != nil {
				log.Printf("Error is %s", err.Error())
			}
			if err := s.deleteUserVerificationToken(ctx, tx, userID); err != nil {
				log.Printf("Error is %s", err.Error())
			}
		} else {
			log.Printf("getUnverifiedUser error is %s", err.Error())
		}

		if err := s.createUser(ctx, tx, user); err != nil {
			return err
		}

		if err := s.createUserVerificationToken(ctx, tx, token, user.ID); err != nil {
			return err
		}

		return nil

	})
}

func (s *UserStore) DeleteUser(ctx context.Context, userID int64) error {
	return util.WithTransaction(s.Db, ctx, func(tx *sql.Tx) error {
		if err := s.deleteUser(ctx, tx, userID); err != nil {
			return err
		}

		if err := s.deleteUserVerificationToken(ctx, tx, userID); err != nil {
			return err
		}

		return nil
	})
}

func (s *UserStore) GetUserByEmail(ctx context.Context, email string) (*User, error) {
	query := `SELECT id, password FROM users WHERE email = $1 AND is_verified = true`

	ctx, cancel := context.WithTimeout(ctx, util.QueryTimeoutDuration)
	defer cancel()

	var user User

	if err := s.Db.QueryRowContext(ctx, query, email).Scan(
		&user.ID,
		&user.Password,
	); err != nil {
		switch err {
		case sql.ErrNoRows:
			return nil, util.ErrorNotFound
		default:
			return nil, err
		}
	}
	return &user, nil
}

func (s *UserStore) getUnverifiedUser(ctx context.Context, tx *sql.Tx, email string) (int64, error) {
	query := `
		SELECT id FROM users WHERE email = $1 AND is_verified = false
	`

	ctx, cancel := context.WithTimeout(ctx, util.QueryTimeoutDuration)
	defer cancel()

	var userID int64

	err := tx.QueryRowContext(ctx, query, email).Scan(
		&userID,
	)

	return userID, err
}

func (s *UserStore) GetUserByID(ctx context.Context, userID int64) (*User, error) {
	query := `SELECT id, username, email FROM users u WHERE u.id = $1`

	ctx, cancel := context.WithTimeout(ctx, util.QueryTimeoutDuration)
	defer cancel()

	var user User
	if err := s.Db.QueryRowContext(ctx, query, userID).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
	); err != nil {
		return nil, err
	}

	return &user, nil
}

func (s *UserStore) CreateRefreshToken(ctx context.Context, userID int64, tokenHash string, expiresAt time.Time) error {
	query := `
		INSERT INTO refresh_tokens (user_id, token_hash, expires_at)
		VALUES ($1, $2, $3)
	`

	ctx, cancel := context.WithTimeout(ctx, util.QueryTimeoutDuration)
	defer cancel()

	_, err := s.Db.ExecContext(ctx, query, userID, tokenHash, expiresAt)
	return err
}

func (s *UserStore) GetUserByRefreshToken(ctx context.Context, tokenHash string) (*User, error) {
	query := `
		SELECT u.id FROM users u INNER JOIN refresh_tokens rt ON u.id = rt.user_id
		WHERE rt.token_hash = $1 AND rt.expires_at > NOW() AND rt.revoked = FALSE
	`

	ctx, cancel := context.WithTimeout(ctx, util.QueryTimeoutDuration)
	defer cancel()

	var user User
	err := s.Db.QueryRowContext(ctx, query, tokenHash).Scan(
		&user.ID,
	)

	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return nil, util.ErrorNotFound
		default:
			return nil, err
		}
	}

	return &user, nil
}

func (s *UserStore) RevokeRefreshToken(ctx context.Context, tokenHash string) error {
	query := `
		UPDATE refresh_tokens SET revoked = TRUE WHERE token_hash = $1
	`
	ctx, cancel := context.WithTimeout(ctx, util.QueryTimeoutDuration)
	defer cancel()

	_, err := s.Db.ExecContext(ctx, query, tokenHash)
	return err
}

func (s *UserStore) DeleteExpiredRefreshTokens(ctx context.Context) error {
	query := `
		DELETE FROM refresh_tokens WHERE expires_at < NOW() OR revoked = TRUE
	`

	ctx, cancel := context.WithTimeout(ctx, util.QueryTimeoutDuration)
	defer cancel()

	_, err := s.Db.ExecContext(ctx, query)
	return err
}

func (s *UserStore) createUser(ctx context.Context, tx *sql.Tx, user *User) error {
	query := `
		INSERT INTO users (email, username, password)
		VALUES ($1, $2, $3) RETURNING id
	`
	ctx, cancel := context.WithTimeout(ctx, util.QueryTimeoutDuration)
	defer cancel()

	err := tx.QueryRowContext(ctx, query, user.Email, user.Username, user.Password).Scan(
		&user.ID,
	)
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

func (s *UserStore) createUserVerificationToken(ctx context.Context, tx *sql.Tx, token string, userID int64) error {
	query := `INSERT INTO users_verification_tracking (token, user_id) VALUES ($1, $2)`

	ctx, cancel := context.WithTimeout(ctx, util.QueryTimeoutDuration)
	defer cancel()

	_, err := tx.ExecContext(ctx, query, token, userID)
	if err != nil {
		return err
	}

	return nil
}

func (s *UserStore) getUserByVerificationToken(ctx context.Context, tx *sql.Tx, token string) (*User, error) {
	query := `
		SELECT u.id FROM users u JOIN users_verification_tracking uv ON u.id = uv.user_id
		WHERE uv.token = $1
	`
	hash := sha256.Sum256([]byte(token))
	hashToken := hex.EncodeToString(hash[:])

	ctx, cancel := context.WithTimeout(ctx, util.QueryTimeoutDuration)
	defer cancel()

	var user User
	if err := tx.QueryRowContext(ctx, query, hashToken).Scan(
		&user.ID,
	); err != nil {
		switch err {
		case sql.ErrNoRows:
			return nil, util.ErrorNotFound
		default:
			return nil, err
		}
	}

	return &user, nil

}

func (s *UserStore) updateUser(ctx context.Context, tx *sql.Tx, user *User) error {
	query := `
		UPDATE users SET is_verified = $1 WHERE id = $2
	`
	user.IsVerified = true
	ctx, cancel := context.WithTimeout(ctx, util.QueryTimeoutDuration)
	defer cancel()

	_, err := tx.ExecContext(ctx, query, user.IsVerified, user.ID)
	if err != nil {
		return err
	}
	return nil
}

func (s *UserStore) deleteUserVerificationToken(ctx context.Context, tx *sql.Tx, userID int64) error {
	query := `DELETE FROM users_verification_tracking WHERE user_id = $1`

	ctx, cancel := context.WithTimeout(ctx, util.QueryTimeoutDuration)
	defer cancel()

	_, err := tx.ExecContext(ctx, query, userID)
	if err != nil {
		return err
	}

	return nil

}

func (s *UserStore) deleteUser(ctx context.Context, tx *sql.Tx, userID int64) error {
	query := `DELETE FROM users WHERE id = $1`

	ctx, cancel := context.WithTimeout(ctx, util.QueryTimeoutDuration)
	defer cancel()

	_, err := tx.ExecContext(ctx, query, userID)
	if err != nil {
		return err
	}

	return nil
}
