package main

import (
	"net/http"

	"github.com/Martins-Iroka/MyGallery-Backend/internal/photo"
	"github.com/Martins-Iroka/MyGallery-Backend/internal/util"
)

type PhotoPostResponsePayload struct {
	ID           int64  `json:"id"`
	Photographer string `json:"photographer"`
	Original     string `json:"original"`
	Large2x      string `json:"large2x"`
	Large        string `json:"large"`
	Medium       string `json:"medium"`
	Small        string `json:"small"`
	Landscape    string `json:"landscape"`
	Tiny         string `json:"tiny"`
}

type CreateCommentRequestPayload struct {
	UserID  int64  `json:"userID" validate:"required"`
	Content string `json:"content" validate:"required"`
}

type PhotoCommentResponsePayload struct {
	Content   string `json:"content" validate:"required"`
	CreatedAt string `json:"created_at" validate:"required"`
	Username  string `json:"username" validate:"required"`
}

func (app *application) createCommentPostHandler(w http.ResponseWriter, r *http.Request) {
	postID := getPostIDFromContext(r)

	var payload CreateCommentRequestPayload
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

	if err := app.store.PicturePost.CreatePhotoComment(ctx, comment); err != nil {
		util.InternalServerErrorResponse(w, r, err, app.logger)
		return
	}

	if err := util.JsonResponse(w, http.StatusCreated, map[string]string{"created": "true"}); err != nil {
		util.InternalServerErrorResponse(w, r, err, app.logger)
	}
}

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

	photoList, err := app.store.PicturePost.GetAllPost(ctx, paginate)
	if err != nil {
		util.InternalServerErrorResponse(w, r, err, app.logger)
		return
	}

	var photos []PhotoPostResponsePayload

	for _, p := range photoList {
		photoPost := &PhotoPostResponsePayload{
			ID:           p.ID,
			Photographer: p.Photographer,
			Original:     p.Original,
			Large2x:      p.Large2x,
			Large:        p.Large,
			Medium:       p.Medium,
			Small:        p.Small,
			Landscape:    p.Landscape,
			Tiny:         p.Tiny,
		}

		photos = append(photos, *photoPost)
	}

	if err := util.JsonResponse(w, http.StatusOK, photos); err != nil {
		util.InternalServerErrorResponse(w, r, err, app.logger)
	}
}

func (app *application) getCommentsByPostID(w http.ResponseWriter, r *http.Request) {
	postID := getPostIDFromContext(r)

	ctx := r.Context()

	comments, err := app.store.PicturePost.GetCommentsByPostID(ctx, postID)
	if err != nil {
		util.InternalServerErrorResponse(w, r, err, app.logger)
		return
	}

	var photoComments []PhotoCommentResponsePayload
	for _, pc := range comments {
		photoComment := &PhotoCommentResponsePayload{
			Content:   pc.Content,
			CreatedAt: pc.CreateAt,
			Username:  pc.Username,
		}

		photoComments = append(photoComments, *photoComment)
	}

	if err := util.JsonResponse(w, http.StatusOK, photoComments); err != nil {
		util.InternalServerErrorResponse(w, r, err, app.logger)
	}
}

func getPostIDFromContext(r *http.Request) int64 {
	postID, _ := r.Context().Value(postCtx).(int64)
	return postID
}
