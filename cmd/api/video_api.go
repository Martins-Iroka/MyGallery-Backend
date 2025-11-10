package main

import (
	"net/http"

	"github.com/Martins-Iroka/MyGallery-Backend/internal/util"
)

type VideoPostResponsePayload struct {
	ID          int64               `json:"id"`
	Video_Url   string              `json:"video_url"`
	Duration    int                 `json:"duration"`
	Video_Files []VideoDownloadFile `json:"video_files"`
}

type VideoDownloadFile struct {
	Video_Link string `json:"video_link"`
	Video_Size int32  `json:"video_size"`
}

func (app *application) getVideosHandler(w http.ResponseWriter, r *http.Request) {

	p := PaginatedFeedQueryAPI{
		Limit:  20,
		Offset: 0,
	}

	p, err := p.parse(r)
	if err != nil {
		util.BadRequestErrorResponse(w, r, err, app.logger)
		return
	}

	if err := util.Validate.Struct(p); err != nil {
		util.BadRequestErrorResponse(w, r, err, app.logger)
		return
	}

	paginate := util.PaginatedPostQuery{
		Limit:  p.Limit,
		Offset: p.Offset,
	}

	ctx := r.Context()

	videoList, err := app.store.VideoPost.GetVideoPostAndDownloadFile(ctx, paginate)
	if err != nil {
		util.InternalServerErrorResponse(w, r, err, app.logger)
		return
	}

	var videos []VideoPostResponsePayload
	for _, v := range videoList {
		video := &VideoPostResponsePayload{
			ID:        v.ID,
			Video_Url: v.Video_Url,
			Duration:  v.Duration,
		}
		var videoDownloadFile []VideoDownloadFile
		for _, vdf := range v.Files {
			df := &VideoDownloadFile{
				vdf.Video_Link,
				vdf.Video_Size,
			}

			videoDownloadFile = append(videoDownloadFile, *df)
		}
		video.Video_Files = videoDownloadFile
		videos = append(videos, *video)
	}

	if err := util.JsonResponse(w, http.StatusOK, videos); err != nil {
		util.InternalServerErrorResponse(w, r, err, app.logger)
	}
}
