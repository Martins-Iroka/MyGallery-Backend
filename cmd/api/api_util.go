package main

type userKey string

const userContextKey userKey = "user"

type photoPostKey string

const photoPostCtx photoPostKey = "photoPostID"

type videoPostKey string

const videoPostCtx videoPostKey = "videoPostID"

// func getUserFromContext(r *http.Request) *user.User {
// 	user, _ := r.Context().Value(userContextKey).(*user.User)
// 	return user
// }
