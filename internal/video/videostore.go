package video

import (
	"context"
	"database/sql"
	"encoding/json"

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
type VideoDownloadFileList []VideoDownloadFile
type VideoPostAndDownloadFile struct {
	VideoPost
	files VideoDownloadFileList
}

type VideoStore struct {
	Db *sql.DB
}

func (v *VideoStore) CreateVideoPost(ctx context.Context, videoPost *VideoPost) error {
	query := `
		INSERT INTO video_posts (id, video_url, duration) VALUES ($1, $2, $3)
	`
	ctx, cancel := context.WithTimeout(ctx, util.QueryTimeoutDuration)
	defer cancel()

	_, err := v.Db.ExecContext(ctx, query, videoPost.ID, videoPost.Video_Url, videoPost.Duration)

	if err != nil {
		return err
	}
	return nil
}

func (v *VideoStore) CreateVideoDownloadFile(ctx context.Context, file *VideoDownloadFile, videoPostId int64) error {
	query := `
		INSERT INTO video_download_files (id, video_post_id, video_link, video_size) VALUES ($1, $2, $3, $4)
	`

	ctx, cancel := context.WithTimeout(ctx, util.QueryTimeoutDuration)
	defer cancel()

	_, err := v.Db.ExecContext(ctx, query, file.ID, videoPostId, file.Video_Link, file.Video_Size)
	if err != nil {
		return err
	}
	return nil
}

func (v *VideoStore) GetVideoPostAndDownloadFile(ctx context.Context, paginate util.PaginatedPostQuery) ([]VideoPostAndDownloadFile, error) {
	query := `
		SELECT 
			vp.id,
			vp.video_url,
			vp.duration,
			COALESCE(
				json_agg(
					json_build_object(
						'id', vdf.id,
						'video_post_id', vdf.video_post_id,
						'video_link', vdf.video_link,
						'video_size', vdf.video_size
					)
				) FILTER (WHERE vdf.id IS NOT NULL),
				'[]'
			) as download_files
		FROM video_posts vp
		LEFT JOIN video_download_files vdf ON vp.id = vdf.video_post_id
		GROUP BY vp.id, vp.video_url, vp.duration
		ORDER BY vp.id DESC
		LIMIT $1 OFFSET $2
	`

	ctx, cancel := context.WithTimeout(ctx, util.QueryTimeoutDuration)
	defer cancel()

	rows, err := v.Db.QueryContext(ctx, query, paginate.Limit, paginate.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []VideoPostAndDownloadFile
	for rows.Next() {
		var vp VideoPostAndDownloadFile
		var filesJSON []byte

		if err := rows.Scan(&vp.ID, &vp.Video_Url, &vp.Duration, &filesJSON); err != nil {
			return nil, err
		}

		// Unmarshal the JSON array of download files
		if err := json.Unmarshal(filesJSON, &vp.files); err != nil {
			return nil, err
		}

		results = append(results, vp)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return results, nil
}
