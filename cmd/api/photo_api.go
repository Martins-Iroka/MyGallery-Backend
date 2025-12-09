package main

import (
	"log"
	"net/http"

	"github.com/Martins-Iroka/MyGallery-Backend/internal/photo"
	"github.com/Martins-Iroka/MyGallery-Backend/internal/util"
)

type PhotoPostResponsePayload struct {
	PhotoItems []PhotoPostItem `json:"photo_items"`
	NextPage   int             `json:"next_page"`
}

type PhotoPostItem struct {
	ID           int64  `json:"id"`
	Photographer string `json:"photographer"`
	Original     string `json:"original"`
	Large2x      string `json:"large2x"`
	Large        string `json:"large"`
	Medium       string `json:"medium"`
	Small        string `json:"small"`
	Portrait     string `json:"portrait"`
	Landscape    string `json:"landscape"`
	Tiny         string `json:"tiny"`
}

type CreatePhotoCommentRequestPayload struct {
	UserID  int64  `json:"userID" validate:"required"`
	Content string `json:"content" validate:"required"`
}

type CreatePhotoCommentResponsePayload struct {
	Created bool `json:"created"`
}

type PhotoCommentResponsePayload struct {
	Content   string `json:"content"`
	CreatedAt string `json:"created_at"`
	Username  string `json:"username"`
	CommentID int64  `json:"id"`
}

// CreateCommentForPhotoPost godoc
//
//	@summary		Create comment for a photo post
//	@Description	Create comment for a photo post
//	@Tags			photos
//	@Accept			json
//	@Produce		json
//	@Param			postID	path		int									true	"Post ID"
//	@Param			payload	body		CreatePhotoCommentRequestPayload	true	"Comment"
//	@Success		201		{object}	CreatePhotoCommentResponsePayload	"Comment created"
//	@Failure		400		{object}	error
//	@Failure		500		{object}	error
//	@Security		ApiKeyAuth
//	@Router			/photos/{postID}/create-comment [post]
func (app *application) createCommentForPostHandler(w http.ResponseWriter, r *http.Request) {
	postID := getPostIDFromContext(r)

	var payload CreatePhotoCommentRequestPayload
	if err := util.ReadJson(w, r, &payload); err != nil {
		util.BadRequestErrorResponse(w, r, err, app.logger)
		return
	}

	if err := util.Validate.Struct(payload); err != nil {
		util.BadRequestErrorResponse(w, r, err, app.logger)
		return
	}

	comment := &photo.PhotoComment{
		Content: payload.Content,
		PostID:  postID,
		UserID:  payload.UserID,
	}

	ctx := r.Context()

	if err := app.store.PhotoPost.CreatePhotoComment(ctx, comment); err != nil {
		util.InternalServerErrorResponse(w, r, err, app.logger)
		return
	}

	if err := util.JsonResponse(w, http.StatusCreated, CreatePhotoCommentResponsePayload{true}); err != nil {
		util.InternalServerErrorResponse(w, r, err, app.logger)
	}
}

// GetPhotoPosts godoc
//
//	@summary		Get all photo posts
//	@Description	Get all photo posts
//	@Tags			photos
//	@Accept			json
//	@Produce		json
//	@Param			limit	query		int	false	"Limit"
//	@Param			offset	query		int	false	"Offset"
//	@Success		200		{object}	PhotoPostResponsePayload
//	@Failure		400		{object}	error
//	@Failure		500		{object}	error
//	@Security		ApiKeyAuth
//	@Router			/photos [get]
func (app *application) getPhotosHandler(w http.ResponseWriter, r *http.Request) {

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

	photoList, err := app.store.PhotoPost.GetAllPost(ctx, paginate)
	if err != nil {
		util.InternalServerErrorResponse(w, r, err, app.logger)
		return
	}

	nextOffset := p.Offset + p.Limit
	log.Print(nextOffset)
	if len(photoList) < p.Limit {
		nextOffset = -1
	}

	photos := []PhotoPostItem{}

	for _, p := range photoList {
		photoPost := &PhotoPostItem{
			ID:           p.ID,
			Photographer: p.Photographer,
			Original:     p.Original,
			Large2x:      p.Large2x,
			Large:        p.Large,
			Medium:       p.Medium,
			Small:        p.Small,
			Portrait:     p.Portrait,
			Landscape:    p.Landscape,
			Tiny:         p.Tiny,
		}

		photos = append(photos, *photoPost)
	}

	photoData := PhotoPostResponsePayload{
		PhotoItems: photos,
		NextPage:   nextOffset,
	}

	if err := util.JsonResponse(w, http.StatusOK, photoData); err != nil {
		util.InternalServerErrorResponse(w, r, err, app.logger)
	}
}

// GetComments godoc
//
//	@summary		Get all comments using a postID
//	@Description	Get all comments using a postID
//	@Tags			photos
//	@Accept			json
//	@Produce		json
//	@Param			postID	path		int	true	"Post ID"
//	@Success		200		{object}	[]PhotoCommentResponsePayload
//	@Failure		400		{object}	error
//	@Failure		500		{object}	error
//	@Security		ApiKeyAuth
//	@Router			/photos/{postID}/comments [get]
func (app *application) getCommentsByPostID(w http.ResponseWriter, r *http.Request) {
	postID := getPostIDFromContext(r)

	ctx := r.Context()

	comments, err := app.store.PhotoPost.GetCommentsByPostID(ctx, postID)
	if err != nil {
		util.InternalServerErrorResponse(w, r, err, app.logger)
		return
	}

	photoComments := []PhotoCommentResponsePayload{}
	for _, pc := range comments {
		photoComment := &PhotoCommentResponsePayload{
			Content:   pc.Content,
			CreatedAt: pc.CreatedAt,
			Username:  pc.Username,
			CommentID: pc.ID,
		}

		photoComments = append(photoComments, *photoComment)
	}

	if err := util.JsonResponse(w, http.StatusOK, photoComments); err != nil {
		util.InternalServerErrorResponse(w, r, err, app.logger)
	}
}

func getPostIDFromContext(r *http.Request) int64 {
	postID, _ := r.Context().Value(photoPostCtx).(int64)
	return postID
}
