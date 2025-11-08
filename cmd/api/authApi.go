package main

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"net/http"
	"time"

	"github.com/Martins-Iroka/MyGallery-Backend/internal/auth"
	"github.com/Martins-Iroka/MyGallery-Backend/internal/user"
	"github.com/Martins-Iroka/MyGallery-Backend/internal/util"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type RegisterUserRequestPayload struct {
	Username string `json:"username" validate:"required,max=100"`
	Email    string `json:"email" validate:"required,email,max=255"`
	Password string `json:"password" validate:"required,min=5,max=72"`
}

type TokenResponsePayload struct {
	Token string `json:"token"`
}

type VerifyUserRequestPayload struct {
	Code  string `json:"code" validate:"required,len=6"`
	Email string `json:"email" validate:"required,email,max=255"`
	Token string `json:"token" validate:"required"`
}

type VerifyUserResponsePayload struct {
	Status string `json:"status"`
}

type LoginUserRequestPayload struct {
	Email    string `json:"email" validate:"required,email,max=255"`
	Password string `json:"password" validate:"required,min=5,max=72"`
}

// RegisterUserHandler godoc
//
//	@summary		Registers a user
//	@Description	Registers a user
//	@Tags			authentication
//	@Accept			json
//	@Produce		json
//	@Param			payload	body		RegisterUserRequestPayload	true	"User credentials"
//	@Success		201		{object}	TokenResponsePayload		"User registered"
//	@Failure		400		{object}	error
//	@Failure		500		{object}	error
//	@Router			/authentication/register [post]
func (app *application) registerUserHandler(w http.ResponseWriter, r *http.Request) {
	var payload RegisterUserRequestPayload

	if err := util.ReadJson(w, r, &payload); err != nil {
		util.BadRequestErrorResponse(w, r, err, app.logger)
		return
	}

	if err := util.Validate.Struct(payload); err != nil {
		util.BadRequestErrorResponse(w, r, err, app.logger)
		return
	}

	hashedPassword, err := auth.HashPassword(payload.Password)
	if err != nil {
		util.InternalServerErrorResponse(w, r, err, app.logger)
		return
	}

	user := &user.User{
		Username: payload.Username,
		Email:    payload.Email,
		Password: hashedPassword,
	}

	ctx := r.Context()

	token := uuid.New().String()

	hash := sha256.Sum256([]byte(token))
	hashToken := hex.EncodeToString(hash[:])

	if err := app.store.User.CreateUserAndVerificationToken(ctx, user, hashToken); err != nil {
		switch err {
		case util.ErrorDuplicateEmail, util.ErrorDuplicateUsername:
			util.BadRequestErrorResponse(w, r, err, app.logger)
		default:
			util.InternalServerErrorResponse(w, r, err, app.logger)
		}
		return
	}

	tokenResponse := TokenResponsePayload{
		Token: token,
	}

	go func(userID int64, email string) {
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
		if err := app.otpVerification.SendVerificationCode(email); err != nil {
			app.logger.Errorw("error sending verification email", "error", err, "user_id", userID)

			if deleteErr := app.store.User.DeleteUser(ctx, userID); deleteErr != nil {
				app.logger.Errorw("error deleting user after email failure",
					"error", deleteErr, "user_id", userID)
			} else {
				app.logger.Infow("user deleted after email failure",
					"user_id", userID)
			}
		}
	}(user.ID, user.Email)

	if err := util.JsonResponse(w, http.StatusCreated, tokenResponse); err != nil {
		util.InternalServerErrorResponse(w, r, err, app.logger)
	}
}

// VerifyUserHandler godoc

// @summary		User verification
// @Description	verify user
// @Tags			authentication
// @Accept			json
// @Produce		json
// @Param			payload	body		VerifyUserRequestPayload	true	"User verification credentials"
// @Success		200		{object}	VerifyUserResponsePayload	"User verified"
// @Failure		400		{object}	error
// @Failure		500		{object}	error
// @Router			/authentication/verify [post]
func (app *application) verifyUserHandler(w http.ResponseWriter, r *http.Request) {
	var payload VerifyUserRequestPayload

	if err := util.ReadJson(w, r, &payload); err != nil {
		util.BadRequestErrorResponse(w, r, err, app.logger)
		return
	}

	if err := util.Validate.Struct(payload); err != nil {
		util.BadRequestErrorResponse(w, r, err, app.logger)
		return
	}

	if err := app.otpVerification.VerifyCode(payload.Email, payload.Code); err != nil {
		util.InternalServerErrorResponse(w, r, err, app.logger)
		return
	}

	if err := app.store.User.ActivateUser(r.Context(), payload.Token); err != nil {
		switch err {
		case util.ErrorNotFound:
			util.NotFoundErrorResponse(w, r, err, app.logger)
		default:
			util.InternalServerErrorResponse(w, r, err, app.logger)
		}
		return
	}

	verifyUserResponse := VerifyUserResponsePayload{Status: "verified"}

	if err := util.JsonResponse(w, http.StatusOK, verifyUserResponse); err != nil {
		util.InternalServerErrorResponse(w, r, err, app.logger)
	}
}

// LoginUserHandler godoc
//
//	@summary	User login
//	@Tags		authentication
//	@Accept		json
//	@Produce	json
//	@Param		payload	body		LoginUserRequestPayload	true	"User login credentials"
//	@Success	200		{string}	Token					"User token"
//	@Failure	400		{object}	error
//	@Failure	500		{object}	error
//	@Router		/authentication/login [post]
func (app *application) loginUserHandler(w http.ResponseWriter, r *http.Request) {
	var payload LoginUserRequestPayload

	if err := util.ReadJson(w, r, &payload); err != nil {
		util.BadRequestErrorResponse(w, r, err, app.logger)
		return
	}

	if err := util.Validate.Struct(payload); err != nil {
		util.BadRequestErrorResponse(w, r, err, app.logger)
		return
	}

	user, err := app.store.User.GetUserByEmail(r.Context(), payload.Email)
	if err != nil {
		switch err {
		case util.ErrorNotFound:
			util.NotFoundErrorResponse(w, r, err, app.logger)
		default:
			util.InternalServerErrorResponse(w, r, err, app.logger)
		}
		return
	}

	if err := auth.ComparePasswords(user.Password, payload.Password); err != nil {
		util.InternalServerErrorResponse(w, r, err, app.logger)
		return
	}

	claims := jwt.MapClaims{
		"sub": user.ID,
		"exp": time.Now().Add(app.config.AuthConfig.Exp).Unix(),
		"iat": time.Now().Unix(),
		"nbf": time.Now().Unix(),
		"iss": app.config.AuthConfig.Iss,
		"aud": app.config.AuthConfig.Iss,
	}

	token, err := app.auth.GenerateToken(claims)
	if err != nil {
		util.InternalServerErrorResponse(w, r, err, app.logger)
		return
	}

	if err := util.JsonResponse(w, http.StatusOK, token); err != nil {
		util.InternalServerErrorResponse(w, r, err, app.logger)
	}
}
