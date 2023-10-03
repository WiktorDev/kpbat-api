package main

import (
	"github.com/foolin/goview"
	"github.com/foolin/goview/supports/echoview-v4"
	"github.com/labstack/echo/v4"
	kpbatApi "kpbatApi/api/base"
	"kpbatApi/api/routing"
)

var config = kpbatApi.LoadConfigFile()

func main() {
	app := echo.New()

	gv := goview.New(goview.Config{
		Root:         "views",                   //template root path
		Extension:    ".html",                   //file extension
		Master:       "layouts/master",          //master layout file
		Partials:     []string{"partials/head"}, //partial files
		DisableCache: true,                      //if disable cache, auto reload template file for debug.
	})

	goview.Use(gv)

	kpbatApi.InitDatabase(config.Database)
	app.Renderer = echoview.Default()

	v1 := app.Group("/v1")
	routing.InitGalleryRouting(v1)
	routing.InitContactRouting(app)
	routing.InitPanelRouting(v1)

	app.Logger.Fatal(app.Start(config.Bind))
}
