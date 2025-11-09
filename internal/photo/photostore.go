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

type PhotoStore struct {
	Db *sql.DB
}

func (p *PhotoStore) CreatePost(ctx context.Context, photo *PhotoPost) error {
	query := `
		INSERT INTO photo_posts (id, photographer, original, large2x, large, medium, small, portrait, landscape, tiny)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
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
		return nil
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
