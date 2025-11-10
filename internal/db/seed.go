package db

import (
	"context"
	"database/sql"
	"log"

	"github.com/Martins-Iroka/MyGallery-Backend/internal"
	"github.com/Martins-Iroka/MyGallery-Backend/internal/photo"
	"github.com/Martins-Iroka/MyGallery-Backend/internal/video"
)

type Seed struct {
	store internal.Storage
}

func NewSeedInstance(db *sql.DB) Seed {
	store := internal.NewStorage(db)
	return Seed{store: store}
}

func (ss Seed) CreatePhotos(photos []photo.PhotoPost) {
	log.Printf("CreatePhotos len %d", len(photos))
	ctx := context.Background()
	for _, photo := range photos {
		err := ss.store.PicturePost.CreatePost(ctx, &photo)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func (ss Seed) CreateVideos(videos []video.VideoPostAndDownloadFile) {
	log.Printf("CreateVideos len %d", len(videos))
	ctx := context.Background()

	for _, v := range videos {
		err := ss.store.VideoPost.CreateVideoPost(ctx, &v.VideoPost)
		if err != nil {
			log.Fatal(err)
		}
		for _, vdf := range v.Files {
			err := ss.store.VideoPost.CreateVideoDownloadFile(ctx, &vdf, vdf.Video_Post_Id)
			if err != nil {
				log.Fatal(err)
			}
		}
	}
}
