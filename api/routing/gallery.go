package routing

import (
	"github.com/labstack/echo/v4"
	kpbatApi "kpbatApi/api/base"
	"kpbatApi/api/base/utils"
	"kpbatApi/api/models"
	"kpbatApi/api/services"
	"net/http"
	"strconv"
)

// @Summary	Find all categories
// @Tags		gallery
// @Accept		json
// @Produce	json
// @Success	200	{object}	models.Category
// @Router		/gallery/categories [get]
func findAllCategories(context echo.Context) error {
	db := kpbatApi.DB()
	var categories []models.Category
	db.Find(&categories)
	return context.JSON(http.StatusOK, categories)
}

// @Summary	Find category by id
// @Tags		gallery
// @Accept		json
// @Produce	json
// @Param		id	path		int	true	"Category ID"
// @Success	200	{object}	models.Category
// @Failure	400	{object}	utils.MessageStruct
// @Failure	500	{object}	utils.MessageStruct
// @Router		/gallery/categories/{id} [get]
func findImages(ctx echo.Context) error {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		return utils.HttpError(ctx, http.StatusBadRequest, utils.Message("Id param must be string!"))
	}
	exists, category := services.FindCategory(id)
	if !exists {
		return utils.HttpError(ctx, http.StatusNotFound, utils.Message("Category not found!"))
	}
	return ctx.JSON(200, category)
}

func InitGalleryRouting(v1 *echo.Group) {
	gallery := v1.Group("/gallery")
	gallery.GET("/categories", findAllCategories)
	gallery.GET("/categories/:id", findImages)
}
