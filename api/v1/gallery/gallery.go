package kpbatApi

import (
	"github.com/labstack/echo/v4"
	kpbatApi "kpbatApi/api/base"
	kpbatApi2 "kpbatApi/api/models"
	"net/http"
)

func findAllCategories(context echo.Context) error {
	db := kpbatApi.DB()
	var categories []kpbatApi2.Category
	db.Find(&categories)
	return context.JSON(http.StatusOK, categories)
}
func createCategory(ctx echo.Context) error {
	bind := new(kpbatApi2.Category)
	db := kpbatApi.DB()

	if err := ctx.Bind(bind); err != nil {
		data := map[string]interface{}{
			"message": err.Error(),
		}
		return ctx.JSON(http.StatusInternalServerError, data)
	}

	if len(bind.DisplayName) < 6 {

	}

	category := &kpbatApi2.Category{
		DisplayName:     bind.DisplayName,
		Description:     bind.Description,
		ImagesDirectory: bind.ImagesDirectory,
	}

	if err := db.Create(&category).Error; err != nil {
		data := map[string]interface{}{
			"message": err.Error(),
		}
		return ctx.JSON(http.StatusInternalServerError, data)
	}

	return ctx.NoContent(http.StatusCreated)
}

func InitRouting(app *echo.Echo) {
	app.GET("/", findAllCategories)
	app.POST("/", createCategory)
}
