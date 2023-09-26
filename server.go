package kpbat_api

import "github.com/labstack/echo/v4"

func main() {
	app := echo.New()
	app.Logger.Fatal(app.Start(":8080"))
}
