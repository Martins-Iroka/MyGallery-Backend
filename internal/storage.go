package internal

import (
	"context"
	"database/sql"
	"time"

	"github.com/Martins-Iroka/MyGallery-Backend/internal/user"
)

type Storage struct {
	User interface {
		ActivateUser(ctx context.Context, token string) error
		CreateUser(ctx context.Context, tx *sql.Tx, user *user.User) error
		CreateAndInviteUser(ctx context.Context, user *user.User, token string, time time.Duration) error
		DeleteUser(ctx context.Context, userID int64) error
	}
}

func NewStorage(db *sql.DB) Storage {
	return Storage{
		User: &user.UserStore{Db: db},
	}
}
