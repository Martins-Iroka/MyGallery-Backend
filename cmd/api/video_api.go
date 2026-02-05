package main

import (
	"net/http"

	"github.com/Martins-Iroka/MyGallery-Backend/internal/util"
	"github.com/Martins-Iroka/MyGallery-Backend/internal/video"
)

type VideoPostInfoResponsePayload struct {
	VideoItems []VideoPostResponsePayload `json:"video_items"`
	NextPage   int                        `json:"next_page"`
}

type VideoPostResponsePayload struct {
	ID          int64               `json:"id"`
	Video_Image string              `json:"video_image"`
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

type CreateVideoCommentResponsePayload struct {
	Created bool `json:"created"`
}

type VideoCommentResponsePayload struct {
	Content   string `json:"content"`
	CreatedAt string `json:"created_at"`
	Username  string `json:"username"`
	ID        int64  `json:"id"`
}

// CreateCommentForVideoPost godoc
//
//	@summary		Create comment for a video post
//	@Description	Create comment for a video post
//	@Tags			videos
//	@Accept			json
//	@Produce		json
//	@Param			postID	path		int															true	"Post ID"
//	@Param			payload	body		CreateVideoCommentRequestPayload							true	"Comment"
//	@Success		201		{object}	util.DataResponse{data=CreateVideoCommentResponsePayload}	"Comment created"
//	@Failure		400		{object}	util.ErrorResponse
//	@Failure		500		{object}	util.ErrorResponse
//	@Security		ApiKeyAuth
//	@Router			/videos/{postID}/create-comment [post]
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

	if err := util.JsonResponse(w, http.StatusCreated, CreateVideoCommentResponsePayload{true}); err != nil {
		util.InternalServerErrorResponse(w, r, err, app.logger)
	}
}

// GetVideoPostComments godoc
//
//	@summary		Get all comments using a postID
//	@Description	Get all comments using a postID
//	@Tags			videos
//	@Accept			json
//	@Produce		json
//	@Param			postID	path		int	true	"Post ID"
//	@Success		200		{object}	util.DataResponse{data=[]VideoCommentResponsePayload}
//	@Failure		400		{object}	util.ErrorResponse
//	@Failure		500		{object}	util.ErrorResponse
//	@Security		ApiKeyAuth
//	@Router			/videos/{postID}/comments [get]
func (app *application) getVideoCommentByPostID(w http.ResponseWriter, r *http.Request) {
	postID := getVideoIdFromContext(r)

	ctx := r.Context()

	comments, err := app.store.VideoPost.GetCommentsByPostID(ctx, postID)
	if err != nil {
		util.InternalServerErrorResponse(w, r, err, app.logger)
		return
	}

	videoComments := []VideoCommentResponsePayload{}
	for _, vc := range comments {
		videoComment := &VideoCommentResponsePayload{
			Content:   vc.Content,
			CreatedAt: vc.CreateAt,
			Username:  vc.Username,
			ID:        vc.ID,
		}
		videoComments = append(videoComments, *videoComment)
	}

	if err := util.JsonResponse(w, http.StatusOK, videoComments); err != nil {
		util.InternalServerErrorResponse(w, r, err, app.logger)
	}
}

// GetVideoPosts godoc
//
//	@summary		Get all video posts
//	@Description	Get all video posts
//	@Tags			videos
//	@Accept			json
//	@Produce		json
//	@Param			limit	query		int	false	"Limit"
//	@Param			offset	query		int	false	"Offset"
//	@Success		200		{object}	util.DataResponse{data=[]VideoPostResponsePayload}
//	@Failure		400		{object}	util.ErrorResponse
//	@Failure		500		{object}	util.ErrorResponse
//	@Security		ApiKeyAuth
//	@Router			/videos [get]
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

	nextOffset := p.Offset + p.Limit
	if len(videoList) < p.Limit {
		nextOffset = -1
	}

	videos := []VideoPostResponsePayload{}
	for _, v := range videoList {
		video := &VideoPostResponsePayload{
			ID:          v.ID,
			Video_Image: v.Video_Image,
			Video_Url:   v.Video_Url,
			Duration:    v.Duration,
		}
		videoDownloadFile := []VideoDownloadFile{}
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

	videoData := VideoPostInfoResponsePayload{
		VideoItems: videos,
		NextPage:   nextOffset,
	}

	if err := util.JsonResponse(w, http.StatusOK, videoData); err != nil {
		util.InternalServerErrorResponse(w, r, err, app.logger)
	}
}

func getVideoIdFromContext(r *http.Request) int64 {
	postID, _ := r.Context().Value(videoPostCtx).(int64)
	return postID
}
