package routing

import (
	"github.com/labstack/echo/v4"
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
	return context.Redirect(http.StatusMovedPermanently, "panel/manage")
}
func renderManage(context echo.Context) error {
	return context.Render(http.StatusOK, "manage", "e1")
}

func InitPanelRouting(v1 *echo.Group) {
	panel := v1.Group("/panel")
	panel.GET("", renderIndex)
	panel.POST("", authorize)
	panel.GET("/manage", renderManage)
}
