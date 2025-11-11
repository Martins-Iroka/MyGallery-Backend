package main

type userKey string

const userContextKey userKey = "user"

type postKey string

const postCtx postKey = "postID"

// func getUserFromContext(r *http.Request) *user.User {
// 	user, _ := r.Context().Value(userContextKey).(*user.User)
// 	return user
// }
