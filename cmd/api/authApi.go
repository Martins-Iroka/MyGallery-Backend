package main

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
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

type RegisterUserResponsePayload struct {
	EmailID string `json:"email_id"`
	Token   string `json:"token"`
}

// RegisterUserHandler godoc
//
//	@summary		Registers a user
//	@Description	Registers a user
//	@Tags			authentication
//	@Accept			json
//	@Produce		json
//	@Param			payload	body		RegisterUserRequestPayload							true	"User credentials"
//	@Success		201		{object}	util.DataResponse{data=RegisterUserResponsePayload}	"User registered"
//	@Failure		400		{object}	util.ErrorResponse
//	@Failure		500		{object}	util.ErrorResponse
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

	if emailId, err := app.otpVerification.SendVerificationCode(user.Email); err != nil {
		app.logger.Errorw("error sending verification email", "error", err, "user_id", user.ID)

		if deleteErr := app.store.User.DeleteUser(ctx, user.ID); deleteErr != nil {
			app.logger.Errorw("error deleting user after email failure",
				"error", deleteErr, "user_id", user.ID)
		} else {
			app.logger.Infow("user deleted after email failure",
				"user_id", user.ID)
		}
	} else {
		tokenResponse := RegisterUserResponsePayload{
			EmailID: emailId,
			Token:   token,
		}
		if err := util.JsonResponse(w, http.StatusCreated, tokenResponse); err != nil {
			util.InternalServerErrorResponse(w, r, err, app.logger)
		}
	}
}

type ResendOTPRequestPayload struct {
	Email string `json:"email" validate:"required,email,max=255"`
}

type ResendOTPResponsePayload struct {
	EmailID string `json:"email_id"`
}

// ResendOTPHandler godoc
//
//	@summary	Resend OTP
//	@Tags		authentication
//	@Accept		json
//	@Produce	json
//	@Param		payload	body		ResendOTPRequestPayload								true	"Resend OTP"
//	@Success	200		{object}	util.DataResponse{data=ResendOTPResponsePayload}	"OTP sent"
//	@Failure	400		{object}	util.ErrorResponse
//	@Failure	500		{object}	util.ErrorResponse
//	@Router		/authentication/resendOTP [post]
func (app *application) resendOTPHandler(w http.ResponseWriter, r *http.Request) {

	var payload ResendOTPRequestPayload

	if err := util.ReadJson(w, r, &payload); err != nil {
		util.BadRequestErrorResponse(w, r, err, app.logger)
		return
	}

	if err := util.Validate.Struct(payload); err != nil {
		util.BadRequestErrorResponse(w, r, err, app.logger)
		return
	}

	if emailId, err := app.otpVerification.SendVerificationCode(payload.Email); err != nil {
		app.logger.Errorw("error resending verification email", "error", err)
		util.InternalServerErrorResponse(w, r, err, app.logger)
	} else {
		response := ResendOTPResponsePayload{
			EmailID: emailId,
		}
		if err := util.JsonResponse(w, http.StatusOK, response); err != nil {
			util.InternalServerErrorResponse(w, r, err, app.logger)
		}
	}

}

type VerifyUserRequestPayload struct {
	Code    string `json:"code" validate:"required,len=6"`
	EmailId string `json:"email_id" validate:"required"`
	Token   string `json:"token" validate:"required"`
}

type VerifyUserResponsePayload struct {
	Status string `json:"status"`
}

// VerifyUserHandler godoc

// @summary		User verification
// @Description	verify user
// @Tags			authentication
// @Accept			json
// @Produce		json
// @Param			payload	body		VerifyUserRequestPayload							true	"User verification credentials"
// @Success		200		{object}	util.DataResponse{data=VerifyUserResponsePayload}	"User verified"
// @Failure		400		{object}	util.ErrorResponse
// @Failure		500		{object}	util.ErrorResponse
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

	if err := app.otpVerification.VerifyCode(payload.EmailId, payload.Code); err != nil {
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

type LoginUserRequestPayload struct {
	Email    string `json:"email" validate:"required,email,max=255"`
	Password string `json:"password" validate:"required,min=5,max=72"`
}

type LoginResponsePayload struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int64  `json:"expires_in"`
	UserID       int64  `json:"user_id"`
}

// LoginUserHandler godoc
//
//	@summary	User login
//	@Tags		authentication
//	@Accept		json
//	@Produce	json
//	@Param		payload	body		LoginUserRequestPayload							true	"User login credentials"
//	@Success	200		{object}	util.DataResponse{data=LoginResponsePayload}	"Login response"
//	@Failure	400		{object}	util.ErrorResponse
//	@Failure	500		{object}	util.ErrorResponse
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
		util.BadRequestErrorResponse(w, r, errors.New("incorrect username or password"), app.logger)
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

	accessToken, err := app.auth.GenerateToken(claims)
	if err != nil {
		util.InternalServerErrorResponse(w, r, err, app.logger)
		return
	}

	refreshToken, err := app.auth.GenerateRefreshToken()
	if err != nil {
		util.InternalServerErrorResponse(w, r, err, app.logger)
		return
	}

	hash := sha256.Sum256([]byte(refreshToken))
	tokenHash := hex.EncodeToString(hash[:])

	refreshExpiry := time.Now().Add(7 * 24 * time.Hour)

	if err := app.store.User.CreateRefreshToken(r.Context(), user.ID, tokenHash, refreshExpiry); err != nil {
		util.InternalServerErrorResponse(w, r, err, app.logger)
		return
	}

	response := LoginResponsePayload{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    int64(app.config.AuthConfig.Exp.Seconds()),
		UserID:       user.ID,
	}

	if err := util.JsonResponse(w, http.StatusOK, response); err != nil {
		util.InternalServerErrorResponse(w, r, err, app.logger)
	}
}

type RefreshTokenRequestPayload struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}

type RefreshTokenResponsePayload struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int64  `json:"expires_in"`
}

// RefreshTokenHandler godoc
//
//	@Summary	Refresh access token
//	@Tags		authentication
//	@Accept		json
//	@Produce	json
//	@Param		payload	body		RefreshTokenRequestPayload							true	"Refresh token"
//	@Success	200		{object}	util.DataResponse{data=RefreshTokenResponsePayload}	"New access token"
//	@Failure	400		{object}	util.ErrorResponse
//	@Failure	401		{object}	util.ErrorResponse
//	@Failure	500		{object}	util.ErrorResponse
//	@Router		/authentication/refresh [post]
func (app *application) refreshTokenHandler(w http.ResponseWriter, r *http.Request) {
	var payload RefreshTokenRequestPayload

	if err := util.ReadJson(w, r, &payload); err != nil {
		util.BadRequestErrorResponse(w, r, err, app.logger)
		return
	}

	if err := util.Validate.Struct(payload); err != nil {
		util.BadRequestErrorResponse(w, r, err, app.logger)
		return
	}

	hash := sha256.Sum256([]byte(payload.RefreshToken))
	tokenHash := hex.EncodeToString(hash[:])

	user, err := app.store.User.GetUserByRefreshToken(r.Context(), tokenHash)
	if err != nil {
		switch err {
		case util.ErrorNotFound:
			util.UnauthorizedErrorResponse(w, r, errors.New("invalid or expired refresh token"), app.logger)
		default:
			util.InternalServerErrorResponse(w, r, err, app.logger)
		}
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

	accessToken, err := app.auth.GenerateToken(claims)
	if err != nil {
		util.InternalServerErrorResponse(w, r, err, app.logger)
		return
	}

	response := RefreshTokenResponsePayload{
		AccessToken: accessToken,
		ExpiresIn:   int64(app.config.AuthConfig.Exp.Seconds()),
	}

	if err := util.JsonResponse(w, http.StatusOK, response); err != nil {
		util.InternalServerErrorResponse(w, r, err, app.logger)
	}
}

type LogoutRequestPayload struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}

// LogoutHandler godoc
//
//	@Summary	Logout user
//	@Tags		authentication
//	@Accept		json
//	@Produce	json
//	@Param		payload	body	LogoutRequestPayload	true	"Refresh token to revoke"
//	@Success	204		"No content"
//	@Failure	400		{object}	util.ErrorResponse
//	@Failure	500		{object}	util.ErrorResponse
//	@Router		/authentication/logout [post]
func (app *application) logoutHandler(w http.ResponseWriter, r *http.Request) {
	var payload LogoutRequestPayload

	if err := util.ReadJson(w, r, &payload); err != nil {
		util.BadRequestErrorResponse(w, r, err, app.logger)
		return
	}

	if err := util.Validate.Struct(payload); err != nil {
		util.BadRequestErrorResponse(w, r, err, app.logger)
		return
	}

	hash := sha256.Sum256([]byte(payload.RefreshToken))
	tokenHash := hex.EncodeToString(hash[:])

	if err := app.store.User.RevokeRefreshToken(r.Context(), tokenHash); err != nil {
		util.InternalServerErrorResponse(w, r, err, app.logger)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
