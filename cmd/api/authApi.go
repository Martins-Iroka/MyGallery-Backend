package main

import (
	"crypto/sha256"
	"encoding/hex"
	"net/http"

	"github.com/Martins-Iroka/MyGallery-Backend/cmd/api/util"
	"github.com/Martins-Iroka/MyGallery-Backend/internal"
	"github.com/Martins-Iroka/MyGallery-Backend/internal/user"
	"github.com/google/uuid"
)

type RegisterUserPayload struct {
	Username string `json:"username" validate:"required,max=100"`
	Email    string `json:"email" validate:"required,max=255"`
	Password string `json:"password" validate:"required,min=5,max=72"`
}

type TokenResponsePayload struct {
	Token string `json:"token"`
}

func (app *application) registerUserHandler(w http.ResponseWriter, r *http.Request) {
	var payload RegisterUserPayload

	if err := util.ReadJson(w, r, &payload); err != nil {
		util.BadRequestErrorResponse(w, r, err, app.logger)
		return
	}

	if err := util.Validate.Struct(payload); err != nil {
		util.BadRequestErrorResponse(w, r, err, app.logger)
		return
	}

	user := &user.User{
		Username: payload.Username,
		Email:    payload.Email,
	}

	ctx := r.Context()

	token := uuid.New().String()

	hash := sha256.Sum256([]byte(token))
	hashToken := hex.EncodeToString(hash[:])

	if err := app.store.User.CreateAndInviteUser(ctx, user, hashToken); err != nil {
		switch err {
		case internal.ErrorDuplicateEmail, internal.ErrorDuplicateUsername:
			util.BadRequestErrorResponse(w, r, err, app.logger)
		default:
			util.InternalServerErrorResponse(w, r, err, app.logger)
		}
		return
	}

	tokenResponse := TokenResponsePayload{
		Token: token,
	}

	if err := app.twilio.SendVerificationCode(user.Email); err != nil {
		app.logger.Errorw("error sending email", "error", err)

		if err := app.store.User.DeleteUser(ctx, user.ID); err != nil {
			app.logger.Errorw("error deleteing user", "error", err)
		}
		util.InternalServerErrorResponse(w, r, err, app.logger)
		return
	}

	if err := util.JsonResponse(w, http.StatusCreated, tokenResponse); err != nil {
		util.InternalServerErrorResponse(w, r, err, app.logger)
	}
}
