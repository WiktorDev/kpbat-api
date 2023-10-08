package routing

import (
	"github.com/labstack/echo/v4"
	kpbatApi "kpbatApi/api/base"
	http_utils "kpbatApi/api/base/utils"
)

func sendMail(ctx echo.Context) error {
	err := kpbatApi.SendMail("wiktor101a@wp.pl", "testowy", "template.html", struct {
		Name    string
		Message string
	}{
		Name:    "Test",
		Message: "wr werw	erreew	erw",
	})
	if err != nil {
		return http_utils.HttpError(ctx, 500, http_utils.Message(err.Error()))
	}
	return ctx.NoContent(200)
}
func InitContactRouting(app *echo.Echo) {
	app.GET("/contact", sendMail)
}
