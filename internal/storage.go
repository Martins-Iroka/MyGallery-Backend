package internal

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/Martins-Iroka/MyGallery-Backend/internal/user"
)

var (
	ErrorNotFound             = errors.New("resource not found")
	ErrorConflict             = errors.New("conflict found modifying resource")
	ErrorUserFollowConflict   = errors.New("you're following this user already")
	ErrorUserUnFollowConflict = errors.New("you're unfollowing this user already")
	ErrorDuplicateEmail       = errors.New("a user with that email already exists")
	ErrorDuplicateUsername    = errors.New("a user with that username already exists")
	QueryTimeoutDuration      = time.Second * 5
)

type Storage struct {
	User interface {
		ActivateUser(ctx context.Context, token string) error
		CreateUser(ctx context.Context, tx *sql.Tx, user *user.User) error
		CreateAndInviteUser(ctx context.Context, user *user.User, token string) error
		DeleteUser(ctx context.Context, userID int64) error
		GetUserByEmail(ctx context.Context, email string) (*user.User, error)
	}
}

func NewStorage(db *sql.DB) Storage {
	return Storage{
		User: &user.UserStore{Db: db},
	}
}
