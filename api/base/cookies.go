package kpbatApi

import (
	"github.com/labstack/echo/v4"
	"net/http"
	"time"
)

func BuildCookie(key string, value string) *http.Cookie {
	cookie := new(http.Cookie)
	cookie.Name = key
	cookie.Value = value
	cookie.Expires = time.Now().Add(24 * time.Hour)
	return cookie
}
func SaveCookie(ctx echo.Context, cookie *http.Cookie) {
	ctx.SetCookie(cookie)
}
func ReadCookie(c echo.Context, key string) string {
	cookie, err := c.Cookie(key)
	if err != nil {
		return ""
	}
	return cookie.Value
}
