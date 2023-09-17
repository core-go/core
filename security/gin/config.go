package gin

import "github.com/gin-gonic/gin"

type SecurityConfig struct {
	SecuritySkip    bool
	Check           func() gin.HandlerFunc
	AuthorizeExact  func(privilege string, action int32) gin.HandlerFunc
	Authorize       func(privilege string, action int32) gin.HandlerFunc
	AuthorizeSub    func(privilege string, sub string, action int32) gin.HandlerFunc
	AuthorizeValues func(values []string) gin.HandlerFunc
}
