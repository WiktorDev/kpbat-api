package services

import (
	"github.com/labstack/echo/v4"
	kpbatApi "kpbatApi/api/base"
	"kpbatApi/api/base/utils"
	"kpbatApi/api/models"
	"net/http"
	"strconv"
)

func ExtractCategory(ctx echo.Context) (*models.Category, error) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		return nil, utils.HttpError(ctx, http.StatusBadRequest, utils.Message("Id param must be string!"))
	}
	exists, category := FindCategory(id)
	if !exists {
		return nil, utils.HttpError(ctx, http.StatusNotFound, utils.Message("Category not found!"))
	}
	return category, nil
}

func FindCategory(id int) (bool, *models.Category) {
	db := kpbatApi.DB()
	var user = models.Category{ID: id}
	err := db.Preload("Images").First(&user).Error
	if err != nil {
		return false, nil
	}
	return true, &user
}
func FindImage(name string) (bool, *models.Image) {
	db := kpbatApi.DB()
	var image = models.Image{FileName: name}
	err := db.First(&image).Error
	if err != nil {
		return false, nil
	}
	return true, &image
}
func RemoveImage(name string) {
	db := kpbatApi.DB()
	_, image := FindImage(name)
	db.Delete(&image)
}
