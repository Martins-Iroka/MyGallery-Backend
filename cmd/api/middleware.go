package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/Martins-Iroka/MyGallery-Backend/internal/user"
	"github.com/Martins-Iroka/MyGallery-Backend/internal/util"
	"github.com/go-chi/chi/v5"
	"github.com/golang-jwt/jwt/v5"
)

func (app *application) authTokenMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			util.UnauthorizedErrorResponse(w, r, errors.New("authorization header is missing"), app.logger)
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			util.UnauthorizedErrorResponse(w, r, errors.New("authorization header is malformed"), app.logger)
			return
		}

		token := parts[1]
		jwtToken, err := app.auth.ValidateToken(token)
		if err != nil {
			util.UnauthorizedErrorResponse(w, r, err, app.logger)
			return
		}

		claims := jwtToken.Claims.(jwt.MapClaims)
		userID, err := strconv.ParseInt(fmt.Sprintf("%.f", claims["sub"]), 10, 64)
		if err != nil {
			util.UnauthorizedErrorResponse(w, r, err, app.logger)
			return
		}

		ctx := r.Context()
		user, err := app.getUserByID(ctx, userID)
		if err != nil {
			util.UnauthorizedErrorResponse(w, r, err, app.logger)
			return
		}

		ctx = context.WithValue(ctx, userContextKey, user)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (app *application) postExistContextMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		idParam := chi.URLParam(r, "postID")
		id, err := strconv.ParseInt(idParam, 10, 64)
		if err != nil {
			util.InternalServerErrorResponse(w, r, err, app.logger)
			return
		}

		ctx := r.Context()

		_, err = app.store.PicturePost.PostExists(ctx, id)
		if err != nil {
			switch err {
			case util.ErrorNotFound:
				util.NotFoundErrorResponse(w, r, err, app.logger)
			default:
				util.InternalServerErrorResponse(w, r, err, app.logger)
			}
			return
		}

		ctx = context.WithValue(ctx, postCtx, id)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (app *application) getUserByID(ctx context.Context, userID int64) (*user.User, error) {
	return app.store.User.GetUserByID(ctx, userID)
}
