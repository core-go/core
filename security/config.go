package security

import "net/http"

type SecurityConfig struct {
	SecuritySkip    bool
	Check           func(next http.Handler) http.Handler
	AuthorizeExact  func(next http.Handler, privilege string, action int32) http.Handler
	Authorize       func(next http.Handler, privilege string, action int32) http.Handler
	AuthorizeSub    func(next http.Handler, privilege string, sub string, action int32) http.Handler
	AuthorizeValues func(next http.Handler, values []string) http.Handler
}
