package main

import (
	"net/http"

	"github.com/Martins-Iroka/MyGallery-Backend/internal/user"
)

type userKey string

const userContextKey userKey = "user"

func getUserFromContext(r *http.Request) *user.User {
	user, _ := r.Context().Value(userContextKey).(*user.User)
	return user
}
