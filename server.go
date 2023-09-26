package main

import (
	"github.com/labstack/echo/v4"
	kpbatApi "kpbatApi/api/base"
	contact "kpbatApi/api/v1/contact"
	gallery "kpbatApi/api/v1/gallery"
)

var config = kpbatApi.LoadConfigFile()

func main() {
	app := echo.New()
	kpbatApi.InitDatabase(config.Database)
	gallery.InitRouting(app)
	contact.InitRouting(app)
	app.Logger.Fatal(app.Start(config.Bind))
}
