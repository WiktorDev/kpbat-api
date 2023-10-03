package routing

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	kpbatApi "kpbatApi/api/base"
	"net/http"
)

func renderIndex(context echo.Context) error {
	return context.Render(http.StatusOK, "index", echo.Map{
		"title": "Authorization",
	})
}
func authorize(context echo.Context) error {
	var token = context.FormValue("token")
	var config = kpbatApi.GetConfig()
	if token != config.Token {
		return context.Render(http.StatusOK, "index", map[string]interface{}{
			"message": "Invalid token!",
			"title":   "Authorization",
		})
	}
	kpbatApi.SaveCookie(context, kpbatApi.BuildCookie("token", token))
	return context.Redirect(http.StatusMovedPermanently, "panel/manage")
}
func renderManage(context echo.Context) error {
	return context.Render(http.StatusOK, "manage", echo.Map{
		"title": "Gallery | Categories",
	})
}

func InitPanelRouting(v1 *echo.Group) {
	var config = kpbatApi.GetConfig()
	panel := v1.Group("/panel")
	panel.Use(middleware.KeyAuthWithConfig(middleware.KeyAuthConfig{
		Skipper: func(c echo.Context) bool {
			return c.Path() == "/v1/panel"
		},
		KeyLookup: "cookie:token",
		Validator: func(token string, c echo.Context) (bool, error) {
			return token == config.Token, nil
		},
	}))
	panel.GET("", renderIndex)
	panel.POST("", authorize)
	panel.GET("/manage", renderManage)
}
