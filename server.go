package main

import (
	"github.com/foolin/goview"
	"github.com/foolin/goview/supports/echoview-v4"
	"github.com/labstack/echo/v4"
	echoSwagger "github.com/swaggo/echo-swagger"
	kpbatApi "kpbatApi/api/base"
	"kpbatApi/api/routing"
	_ "kpbatApi/docs"
)

var config = kpbatApi.LoadConfigFile()

//	@title			Kpbat API
//	@version		1.0
//	@description	Simple REST API for kpbat.com website

//	@license.name	Apache 2.0
//	@license.url	http://www.apache.org/licenses/LICENSE-2.0.html

// @host		api.kpbat.com
// @BasePath	/v1
func main() {
	app := echo.New()
	app.GET("/swagger/*", echoSwagger.WrapHandler)
	app.Static("/resources", "resources")
	gv := goview.New(goview.Config{
		Root:         "views",                   //template root path
		Extension:    ".html",                   //file extension
		Master:       "layouts/master",          //master layout file
		Partials:     []string{"partials/head"}, //partial files
		DisableCache: false,                     //if disable cache, auto reload template file for debug.
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
