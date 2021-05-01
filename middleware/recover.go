package middleware

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

// PanicHandler refer to https://medium.com/@masnun/panic-recovery-middleware-for-go-http-handlers-51147c941f9 and  http://www.golangtraining.in/lessons/middleware/recovering-from-panic.html
func PanicHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		panic("forcing a panic")
	})
}
func Recover(log func(ctx context.Context, msg string)) func(h http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				er := recover()
				if er != nil {
					s := GetError(er)
					log(r.Context(), s)
					w.Header().Set("Content-Type", "application/json")
					w.WriteHeader(http.StatusInternalServerError)
					jsonBody, _ := json.Marshal(map[string]string{
						"error": "Internal Server Error",
					})
					w.Write(jsonBody)
				}
			}()
			h.ServeHTTP(w, r)
		}
		return http.HandlerFunc(fn)
	}
}

func GetError(er interface{}) string {
	if er == nil {
		return "Internal Server Error"
	}
	switch x := er.(type) {
	case string:
		return er.(string)
	case error:
		err := x
		return err.Error()
	default:
		return fmt.Sprintf("%v", er)
	}
}
