package routing

import (
	"fmt"
	"github.com/labstack/echo/v4"
	kpbatApi "kpbatApi/api/base"
	"kpbatApi/api/base/utils"
	"kpbatApi/api/models"
	"net/http"
	"strconv"
)

func findAllCategories(context echo.Context) error {
	db := kpbatApi.DB()
	var categories []models.Category
	db.Find(&categories)
	return context.JSON(http.StatusOK, categories)
}
func createCategory(ctx echo.Context) error {
	db := kpbatApi.DB()
	bind := new(models.Category)

	if err := ctx.Bind(bind); err != nil {
		return utils.HttpError(ctx, http.StatusInternalServerError, utils.Message(err.Error()))
	}
	if err := utils.Validate(ctx, models.CategoryValidator(bind)); err != nil {
		return err
	}

	if err := utils.CreateDirectory(bind.ImagesDirectory); err != nil {
		return utils.HttpError(ctx, http.StatusBadRequest, utils.Message(err.Error()))
	}

	if err := db.Create(&models.Category{
		DisplayName:     bind.DisplayName,
		Description:     bind.Description,
		ImagesDirectory: bind.ImagesDirectory,
	}).Error; err != nil {
		return utils.HttpError(ctx, http.StatusInternalServerError, utils.Message(err.Error()))
	}

	return ctx.NoContent(http.StatusCreated)
}
func deleteCategory(ctx echo.Context) error {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		return utils.HttpError(ctx, http.StatusBadRequest, utils.Message("Id param must be string!"))
	}
	isExists(id)
	//db := kpbatApi.DB()
	//db.Delete(&models.Category{}, id)
	return ctx.NoContent(http.StatusOK)
}

func isExists(id int) {
	db := kpbatApi.DB()

	fmt.Println(db.First(models.Category{}, id))
}

func InitGalleryRouting(v1 *echo.Group) {
	gallery := v1.Group("/gallery")
	gallery.GET("/categories", findAllCategories)
	gallery.POST("/categories", createCategory)
	gallery.DELETE("/categories/:id", deleteCategory)
}
