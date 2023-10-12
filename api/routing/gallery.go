package routing

import (
	"github.com/labstack/echo/v4"
	kpbatApi "kpbatApi/api/base"
	"kpbatApi/api/base/utils"
	"kpbatApi/api/models"
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

// @Summary	Create new categories
// @Tags		gallery
// @Accept		json
// @Produce	json
//
// @Security	ApiKeyAuth
// @param		Authorization	header	string	true	"Authorization"
//
// @Success	201
// @Failure	400	{object}	utils.MessageStruct
// @Failure	500	{object}	utils.MessageStruct
// @Router		/gallery/categories [post]
func createCategory(ctx echo.Context) error {
	db := kpbatApi.DB()
	bind := new(models.Category)

	if err := ctx.Bind(bind); err != nil {
		return utils.HttpError(ctx, http.StatusInternalServerError, utils.Message(err.Error()))
	}
	err, isValid := utils.Validate(ctx, models.CategoryValidator(bind))
	if err != nil {
		return err
	}
	if !isValid {
		return ctx.NoContent(http.StatusBadRequest)
	}

	category := &models.Category{
		DisplayName: bind.DisplayName,
		Description: bind.Description,
	}

	if err := db.Create(&category).Error; err != nil {
		return utils.HttpError(ctx, http.StatusInternalServerError, utils.Message(err.Error()))
	}
	if err := utils.CreateDirectory("category_" + strconv.Itoa(category.ID)); err != nil {
		return utils.HttpError(ctx, http.StatusInternalServerError, utils.Message(err.Error()))
	}
	return ctx.NoContent(http.StatusCreated)
}

// @Summary	Delete category
// @Tags		gallery
// @Accept		json
// @Produce	json
// @Param		id	path	int	true	"Category ID"
// @Success	200
// @Failure	500	{object}	utils.MessageStruct
// @Router		/gallery/categories/{id} [delete]
func deleteCategory(ctx echo.Context) error {
	return ctx.NoContent(http.StatusOK)
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
	exists, category := findCategory(id)
	if !exists {
		return utils.HttpError(ctx, http.StatusNotFound, utils.Message("Category not found!"))
	}
	return ctx.JSON(200, category)
}

func findCategory(id int) (bool, *models.Category) {
	db := kpbatApi.DB()
	var user = models.Category{ID: id}
	err := db.Preload("Images").First(&user).Error
	if err != nil {
		return false, nil
	}
	return true, &user
}
func findCategoryContext(ctx echo.Context) (*models.Category, error) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		return nil, utils.HttpError(ctx, http.StatusBadRequest, utils.Message("Id param must be string!"))
	}
	exists, category := findCategory(id)

	if !exists {
		return nil, utils.HttpError(ctx, http.StatusNotFound, utils.Message("Category not found!"))
	}
	return category, nil
}
func InitGalleryRouting(v1 *echo.Group) {
	gallery := v1.Group("/gallery")
	gallery.GET("/categories", findAllCategories)
	gallery.GET("/categories/:id", findImages)
	gallery.POST("/categories", createCategory)
	gallery.DELETE("/categories/:id", deleteCategory)
}
