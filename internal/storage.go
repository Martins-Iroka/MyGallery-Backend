package internal

import (
	"context"
	"database/sql"

	"github.com/Martins-Iroka/MyGallery-Backend/internal/user"
)

type Storage struct {
	User interface {
		ActivateUser(ctx context.Context, token string) error
		CreateUserAndVerificationToken(ctx context.Context, user *user.User, token string) error
		DeleteUser(ctx context.Context, userID int64) error
		GetUserByEmail(ctx context.Context, email string) (*user.User, error)
	}
}

func NewStorage(db *sql.DB) Storage {
	return Storage{
		User: &user.UserStore{Db: db},
	}
}
