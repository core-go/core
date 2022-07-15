package echo

import (
	"github.com/core-go/core/health"
	"github.com/labstack/echo/v4"
	"net/http"
)

type Handler struct {
	Checkers []health.Checker
}

func NewHandler(checkers ...health.Checker) *Handler {
	return &Handler{checkers}
}

func (c *Handler) Check(ctx echo.Context) error {
	result := health.Check(ctx.Request().Context(), c.Checkers)
	if result.Status == health.StatusUp {
		return ctx.JSON(http.StatusOK, result)
	} else {
		return ctx.JSON(http.StatusInternalServerError, result)
	}
}
