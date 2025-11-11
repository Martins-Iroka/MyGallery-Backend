package main

import (
	"net/http"

	"github.com/Martins-Iroka/MyGallery-Backend/internal/util"
	"github.com/Martins-Iroka/MyGallery-Backend/internal/video"
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

type CreateVideoCommentRequestPayload struct {
	UserID  int64  `json:"userID" validate:"required"`
	Content string `json:"content" validate:"required"`
}

type VideoCommentResponsePayload struct {
	Content   string `json:"content" validate:"required"`
	CreatedAt string `json:"created_at" validate:"required"`
	Username  string `json:"username" validate:"required"`
}

func (app *application) createVideoCommentHandler(w http.ResponseWriter, r *http.Request) {
	postID := getVideoIdFromContext(r)

	var payload CreateVideoCommentRequestPayload
	if err := util.ReadJson(w, r, &payload); err != nil {
		util.BadRequestErrorResponse(w, r, err, app.logger)
		return
	}

	if err := util.Validate.Struct(payload); err != nil {
		util.BadRequestErrorResponse(w, r, err, app.logger)
		return
	}

	comment := &video.VideoComment{
		Content: payload.Content,
		PostID:  postID,
		UserID:  payload.UserID,
	}

	ctx := r.Context()

	if err := app.store.VideoPost.CreateVideoComment(ctx, comment); err != nil {
		util.InternalServerErrorResponse(w, r, err, app.logger)
		return
	}

	if err := util.JsonResponse(w, http.StatusCreated, map[string]string{"created": "true"}); err != nil {
		util.InternalServerErrorResponse(w, r, err, app.logger)
	}
}

func (app *application) getVideoCommentByPostID(w http.ResponseWriter, r *http.Request) {
	postID := getVideoIdFromContext(r)

	ctx := r.Context()

	comments, err := app.store.VideoPost.GetCommentsByPostID(ctx, postID)
	if err != nil {
		util.InternalServerErrorResponse(w, r, err, app.logger)
		return
	}

	var videoComments []VideoCommentResponsePayload
	for _, vc := range comments {
		videoComment := &VideoCommentResponsePayload{
			Content:   vc.Content,
			CreatedAt: vc.CreateAt,
			Username:  vc.Username,
		}
		videoComments = append(videoComments, *videoComment)
	}

	if err := util.JsonResponse(w, http.StatusOK, videoComments); err != nil {
		util.InternalServerErrorResponse(w, r, err, app.logger)
	}
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

func getVideoIdFromContext(r *http.Request) int64 {
	postID, _ := r.Context().Value(videoPostCtx).(int64)
	return postID
}
