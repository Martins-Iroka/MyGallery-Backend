package video

import (
	"context"
	"database/sql"
	"encoding/json"

	"github.com/Martins-Iroka/MyGallery-Backend/internal/util"
)

type VideoPost struct {
	ID          int64
	Video_Image string
	Video_Url   string
	Duration    int
}

type VideoDownloadFile struct {
	ID            int64
	Video_Post_Id int64
	Video_Link    string
	Video_Size    int32
}

type VideoComment struct {
	ID       int64
	PostID   int64
	UserID   int64
	Content  string
	CreateAt string
	Username string
}

type VideoDownloadFileList []VideoDownloadFile
type VideoPostAndDownloadFile struct {
	VideoPost
	Files VideoDownloadFileList
}

type VideoStore struct {
	Db *sql.DB
}

func (v *VideoStore) CreateVideoPost(ctx context.Context, videoPost *VideoPost) error {
	query := `
		INSERT INTO video_posts (id, video_image, video_url, duration) VALUES ($1, $2, $3, $4)
		ON CONFLICT (id) DO UPDATE SET id = EXCLUDED.id, video_image = EXCLUDED.video_image, video_url = EXCLUDED.video_url,
    	duration = EXCLUDED.duration
	`
	ctx, cancel := context.WithTimeout(ctx, util.QueryTimeoutDuration)
	defer cancel()

	_, err := v.Db.ExecContext(ctx, query, videoPost.ID, videoPost.Video_Image, videoPost.Video_Url, videoPost.Duration)

	if err != nil {
		return err
	}
	return nil
}

func (v *VideoStore) CreateVideoDownloadFile(ctx context.Context, file *VideoDownloadFile, videoPostId int64) error {
	query := `
		INSERT INTO video_download_files (id, video_post_id, video_link, video_size) VALUES ($1, $2, $3, $4)
		ON CONFLICT (id) DO NOTHING
	`

	ctx, cancel := context.WithTimeout(ctx, util.QueryTimeoutDuration)
	defer cancel()

	_, err := v.Db.ExecContext(ctx, query, file.ID, videoPostId, file.Video_Link, file.Video_Size)
	if err != nil {
		return err
	}
	return nil
}

func (v *VideoStore) CreateVideoComment(ctx context.Context, video *VideoComment) error {
	query := `INSERT INTO video_comments (post_id, user_id, content) VALUES ($1, $2, $3)`

	ctx, cancel := context.WithTimeout(ctx, util.QueryTimeoutDuration)
	defer cancel()

	_, err := v.Db.ExecContext(ctx, query, video.PostID, video.UserID, video.Content)
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
			vp.video_image,
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

		if err := rows.Scan(&vp.ID, &vp.Video_Url, &vp.Video_Image, &vp.Duration, &filesJSON); err != nil {
			return nil, err
		}

		// Unmarshal the JSON array of download files
		if err := json.Unmarshal(filesJSON, &vp.Files); err != nil {
			return nil, err
		}

		results = append(results, vp)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return results, nil
}

func (v *VideoStore) PostExists(ctx context.Context, videoID int64) error {
	query := `SELECT EXISTS(SELECT 1 FROM video_posts WHERE id = $1)`

	_, err := v.Db.ExecContext(ctx, query, videoID)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return util.ErrorNotFound
		default:
			return err
		}
	}

	return nil
}

func (v *VideoStore) GetCommentsByPostID(ctx context.Context, videoID int64) ([]VideoComment, error) {
	query := `
		SELECT vc.content, vc.created_at, vc.id, u.username FROM video_comments vc JOIN users u ON u.id = vc.user_id
		WHERE vc.post_id = $1 ORDER BY vc.created_at DESC
	`

	ctx, cancel := context.WithTimeout(ctx, util.QueryTimeoutDuration)
	defer cancel()

	rows, err := v.Db.QueryContext(ctx, query, videoID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	comments := []VideoComment{}
	for rows.Next() {
		var vc VideoComment
		err := rows.Scan(
			&vc.Content,
			&vc.CreateAt,
			&vc.ID,
			&vc.Username,
		)
		if err != nil {
			return nil, err
		}
		comments = append(comments, vc)
	}
	return comments, nil
}
