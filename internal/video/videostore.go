package video

import (
	"context"
	"database/sql"

	"github.com/Martins-Iroka/MyGallery-Backend/internal/util"
)

type VideoPost struct {
	ID        int64
	Video_Url string
	Duration  int
}

type VideoDownloadFile struct {
	ID            int64
	Video_Post_Id int64
	Video_Link    string
	Video_Size    int32
}

type VideoPostAndDownloadFile struct {
	VideoPost
	files []VideoDownloadFile
}

type VideoStore struct {
	Db *sql.DB
}

func (v *VideoStore) CreateVideoPost(ctx context.Context, videoPost *VideoPost) error {

	return nil
}

func (v *VideoStore) CreateVideoDownloadFile(ctx context.Context, file *VideoDownloadFile, videoPostId int64) error {

	return nil
}

func (v *VideoStore) GetVideoPostAndDownloadFile(ctx context.Context, paginate util.PaginatedPostQuery) ([]VideoPostAndDownloadFile, error) {

	return nil, nil
}
