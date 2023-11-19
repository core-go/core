package core

import (
	"context"
	"net/http"
)

var Id = "id"
var UserID = "userID"
var UserId = "userId"
func ApplyUserId(str string) {
	UserId = str
}
func GetUser(ctx context.Context, opt...string) (string, bool) {
	user := UserId
	if len(opt) > 0 && len(opt[0]) > 0 {
		user = opt[0]
	}
	u := ctx.Value(user)
	if u != nil {
		u2, ok2 := u.(string)
		if ok2 {
			return u2, ok2
		}
	}
	return "", false
}
func GetString(ctx context.Context, key string) (string, bool) {
	u := ctx.Value(key)
	if u != nil {
		u2, ok2 := u.(string)
		if ok2 {
			return u2, true
		}
	}
	return "", false
}
func RequireUser(ctx context.Context, w http.ResponseWriter, opt...string) (string, bool) {
	userId, ok := GetUser(ctx, opt...)
	if ok {
		return userId, ok
	} else {
		http.Error(w, "cannot get current user", http.StatusForbidden)
		return "", false
	}
}
