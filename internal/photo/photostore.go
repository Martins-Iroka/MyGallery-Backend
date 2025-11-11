package photo

import (
	"context"
	"database/sql"

	"github.com/Martins-Iroka/MyGallery-Backend/internal/util"
)

type PhotoPost struct {
	ID           int64
	Photographer string
	Original     string
	Large2x      string
	Large        string
	Medium       string
	Small        string
	Portrait     string
	Landscape    string
	Tiny         string
}

type PhotoComment struct {
	ID       int64
	PostID   int64
	UserID   int64
	Content  string
	CreateAt string
	Username string
}

type PhotoStore struct {
	Db *sql.DB
}

func (p *PhotoStore) CreatePost(ctx context.Context, photo *PhotoPost) error {
	query := `
		INSERT INTO photo_posts (id, photographer, original, large2x, large, medium, small, portrait, landscape, tiny)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10) ON CONFLICT (id) DO NOTHING
	`

	ctx, cancel := context.WithTimeout(ctx, util.QueryTimeoutDuration)
	defer cancel()

	_, err := p.Db.ExecContext(
		ctx,
		query,
		photo.ID, photo.Photographer, photo.Original, photo.Large2x, photo.Large, photo.Medium, photo.Small,
		photo.Portrait, photo.Landscape, photo.Tiny,
	)

	if err != nil {
		return err
	}
	return nil
}

func (p *PhotoStore) GetAllPost(ctx context.Context, pagination util.PaginatedPostQuery) ([]PhotoPost, error) {
	query := `SELECT * FROM photo_posts p ORDER BY p.photographer LIMIT $1 OFFSET $2`

	ctx, cancel := context.WithTimeout(ctx, util.QueryTimeoutDuration)
	defer cancel()

	rows, err := p.Db.QueryContext(ctx, query, pagination.Limit, pagination.Offset)
	if err != nil {
		return nil, err
	}

	var photos []PhotoPost

	for rows.Next() {
		var photo PhotoPost
		err := rows.Scan(
			&photo.ID,
			&photo.Photographer,
			&photo.Original,
			&photo.Large2x,
			&photo.Large,
			&photo.Medium,
			&photo.Small,
			&photo.Portrait,
			&photo.Landscape,
			&photo.Tiny,
		)
		if err != nil {
			return nil, err
		}

		photos = append(photos, photo)
	}

	return photos, nil
}

func (p *PhotoStore) CreatePhotoComment(ctx context.Context, comment *PhotoComment) error {
	query := `INSERT INTO photos_comment (post_id, user_id, content) VALUES($1, $2, $3)`

	ctx, cancel := context.WithTimeout(ctx, util.QueryTimeoutDuration)
	defer cancel()

	_, err := p.Db.ExecContext(ctx, query, comment.PostID, comment.UserID, comment.Content)
	if err != nil {
		return err
	}

	return nil
}

func (p *PhotoStore) PostExists(ctx context.Context, postID int64) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM photo_posts WHERE id = $1)`

	var exists bool
	err := p.Db.QueryRowContext(ctx, query, postID).Scan(&exists)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return false, util.ErrorNotFound
		default:
			return false, err
		}
	}
	return exists, nil
}

func (p *PhotoStore) GetCommentsByPostID(ctx context.Context, postID int64) ([]PhotoComment, error) {
	query := `
		SELECT pc.content, pc.created_at, u.username FROM photos_comment pc JOIN users u ON u.id = pc.user_id
		WHERE pc.post_id = $1 ORDER BY pc.created_at DESC
	`

	ctx, cancel := context.WithTimeout(ctx, util.QueryTimeoutDuration)
	defer cancel()

	rows, err := p.Db.QueryContext(ctx, query, postID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	comments := []PhotoComment{}
	for rows.Next() {
		var pc PhotoComment
		err := rows.Scan(
			&pc.Content,
			&pc.CreateAt,
			&pc.Username,
		)
		if err != nil {
			return nil, err
		}
		comments = append(comments, pc)
	}

	return comments, nil
}
