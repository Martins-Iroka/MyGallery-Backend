package internal

import (
	"context"
	"database/sql"

	"github.com/Martins-Iroka/MyGallery-Backend/internal/photo"
	"github.com/Martins-Iroka/MyGallery-Backend/internal/user"
	"github.com/Martins-Iroka/MyGallery-Backend/internal/util"
	"github.com/Martins-Iroka/MyGallery-Backend/internal/video"
)

type Storage struct {
	User interface {
		ActivateUser(ctx context.Context, token string) error
		CreateUserAndVerificationToken(ctx context.Context, user *user.User, token string) error
		DeleteUser(ctx context.Context, userID int64) error
		GetUserByEmail(ctx context.Context, email string) (*user.User, error)
		GetUserByID(ctx context.Context, userID int64) (*user.User, error)
	}

	PhotoPost interface {
		CreatePost(context.Context, *photo.PhotoPost) error
		GetAllPost(context.Context, util.PaginatedPostQuery) ([]photo.PhotoPost, error)
		CreatePhotoComment(context.Context, *photo.PhotoComment) error
		PostExists(context.Context, int64) (bool, error)
		GetCommentsByPostID(context.Context, int64) ([]photo.PhotoComment, error)
	}
	VideoPost interface {
		CreateVideoPost(context.Context, *video.VideoPost) error
		CreateVideoDownloadFile(context.Context, *video.VideoDownloadFile, int64) error
		GetVideoPostAndDownloadFile(context.Context, util.PaginatedPostQuery) ([]video.VideoPostAndDownloadFile, error)
		CreateVideoComment(context.Context, *video.VideoComment) error
		PostExists(context.Context, int64) error
		GetCommentsByPostID(context.Context, int64) ([]video.VideoComment, error)
	}
}

func NewStorage(db *sql.DB) Storage {
	return Storage{
		User:      &user.UserStore{Db: db},
		PhotoPost: &photo.PhotoStore{Db: db},
		VideoPost: &video.VideoStore{Db: db},
	}
}

func NewMockStorate() Storage {
	return Storage{
		User: &user.MockUserStore{},
	}
}
