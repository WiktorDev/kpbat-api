package utils

import (
	"github.com/labstack/echo/v4"
	"net/http"
)

func HttpError(ctx echo.Context, code int, i interface{}) error {
	return ctx.JSON(code, i)
}
func Message(text string) interface{} {
	return map[string]interface{}{
		"message": text,
	}
}
func Validate(ctx echo.Context, validated Validated) error {
	if !validated.Ok {
		return HttpError(ctx, http.StatusBadRequest, Message(validated.Message))
	}
	return nil
}

type Validated struct {
	Ok      bool
	Message string
}