package util

import (
	"context"
	"database/sql"
	"errors"
	"time"
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

type PaginatedPostQuery struct {
	Limit  int
	Offset int
}

func WithTransaction(db *sql.DB, ctx context.Context, fn func(*sql.Tx) error) error {
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	if err := fn(tx); err != nil {
		_ = tx.Rollback()
		return err
	}

	return tx.Commit()
}
