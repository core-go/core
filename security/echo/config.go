package echo

import "github.com/labstack/echo/v4"

type SecurityConfig struct {
	SecuritySkip    bool
	Check           func() echo.MiddlewareFunc
	AuthorizeExact  func(privilege string, action int32) echo.MiddlewareFunc
	Authorize       func(privilege string, action int32) echo.MiddlewareFunc
	AuthorizeSub    func(privilege string, sub string, action int32) echo.MiddlewareFunc
	AuthorizeValues func(values []string) echo.MiddlewareFunc
}
