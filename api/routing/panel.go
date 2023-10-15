package routing

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"io"
	kpbatApi "kpbatApi/api/base"
	"kpbatApi/api/base/utils"
	"kpbatApi/api/models"
	"kpbatApi/api/services"
	"net/http"
	"os"
	"path/filepath"
	"slices"
	"strconv"
)

func renderIndex(context echo.Context) error {
	return context.Render(http.StatusOK, "index", echo.Map{
		"title": "Authorization",
	})
}
func authorize(context echo.Context) error {
	var token = context.FormValue("token")
	var config = kpbatApi.GetConfig()
	if token != config.Token {
		return context.Render(http.StatusOK, "index", map[string]interface{}{
			"message": "Invalid token!",
			"title":   "Authorization",
		})
	}
	kpbatApi.SaveCookie(context, kpbatApi.BuildCookie("token", token))
	return context.Redirect(http.StatusMovedPermanently, "panel/manage")
}
func renderManage(context echo.Context) error {
	return context.Render(http.StatusOK, "manage", echo.Map{
		"title": "Gallery | Categories",
	})
}
func createCategory(ctx echo.Context) error {
	db := kpbatApi.DB()
	bind := models.Category{
		DisplayName: ctx.FormValue("name"),
		Description: ctx.FormValue("description"),
	}
	validated := models.CategoryValidator(&bind)
	if !validated.Ok {
		return ctx.Render(http.StatusBadRequest, "manage", echo.Map{
			"title":   "Gallery | Categories",
			"message": validated.Message,
		})
	}
	if err := db.Create(&bind).Error; err != nil {
		return utils.HttpError(ctx, http.StatusInternalServerError, utils.Message(err.Error()))
	}
	if err := utils.CreateDirectory("category_" + strconv.Itoa(bind.ID)); err != nil {
		return utils.HttpError(ctx, http.StatusInternalServerError, utils.Message(err.Error()))
	}
	return ctx.Render(http.StatusCreated, "manage", echo.Map{
		"title": "Gallery | Categories",
	})
}
func renderManageCategory(ctx echo.Context) error {
	return ctx.Render(http.StatusOK, "images", echo.Map{
		"title":      "Gallery | Images",
		"categoryId": ctx.Param("id"),
	})
}
func uploadImages(ctx echo.Context) error {
	db := kpbatApi.DB()

	category, categoryErr := services.ExtractCategory(ctx)
	if categoryErr != nil {
		return categoryErr
	}

	form, err := ctx.MultipartForm()
	if err != nil {
		return err
	}
	files := form.File["files"]

	for _, file := range files {
		availableMimeTypes := []string{"image/png", "image/jpeg"}
		if !slices.Contains(availableMimeTypes, file.Header.Get("Content-Type")) {
			continue
		}

		src, err := file.Open()
		if err != nil {
			return err
		}
		defer src.Close()

		newFileName := uuid.New().String() + filepath.Ext(file.Filename)
		dst, err := os.Create("resources/category_" + strconv.Itoa(category.ID) + "/" + newFileName)
		if err != nil {
			return err
		}
		defer dst.Close()

		if _, err = io.Copy(dst, src); err != nil {
			return err
		}
		if err := db.Create(&models.Image{
			CategoryID: category.ID,
			FileName:   newFileName,
		}).Error; err != nil {
			return utils.HttpError(ctx, http.StatusInternalServerError, utils.Message(err.Error()))
		}
	}

	return ctx.HTML(http.StatusOK, fmt.Sprintf("<p>Uploaded successfully %d files</p>", len(files)))
}

func InitPanelRouting(v1 *echo.Group) {
	var config = kpbatApi.GetConfig()
	panel := v1.Group("/panel")
	panel.Use(middleware.KeyAuthWithConfig(middleware.KeyAuthConfig{
		Skipper: func(c echo.Context) bool {
			return c.Path() == "/v1/panel"
		},
		KeyLookup: "cookie:token",
		Validator: func(token string, c echo.Context) (bool, error) {
			return token == config.Token, nil
		},
	}))
	panel.GET("", renderIndex)
	panel.POST("", authorize)
	panel.GET("/manage", renderManage)
	panel.POST("/manage", createCategory)
	panel.GET("/manage/:id", renderManageCategory)
	panel.POST("/manage/:id", uploadImages)
}
