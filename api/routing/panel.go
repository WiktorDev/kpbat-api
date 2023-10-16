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

func view(template string, title string, context echo.Context, data echo.Map) error {
	if data == nil {
		data = echo.Map{}
	}
	data["title"] = title
	return context.Render(http.StatusOK, template, data)
}

func renderIndex(context echo.Context) error {
	return view("index", "Authorization", context, nil)
}
func authorize(context echo.Context) error {
	var token = context.FormValue("token")
	var config = kpbatApi.GetConfig()
	if token != config.Token {
		return view("index", "Authorization", context, echo.Map{
			"message": "Invalid token!",
		})
	}
	kpbatApi.SaveCookie(context, kpbatApi.BuildCookie("token", token))
	return context.Redirect(http.StatusMovedPermanently, "panel/manage")
}
func renderManage(context echo.Context) error {
	return view("manage", "Gallery | Categories", context, nil)
}
func removeCategory(ctx echo.Context) error {
	db := kpbatApi.DB()
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		return utils.HttpError(ctx, http.StatusBadRequest, utils.Message("Id param must be integer!"))
	}
	_, category := services.FindCategory(id)
	if err := db.Delete(&category).Association("Images").Clear(); err != nil {
		fmt.Println(err)
	}
	utils.RemoveDir(fmt.Sprintf("category_%d", id))
	return ctx.Redirect(http.StatusMovedPermanently, "/v1/panel/manage")
}
func createCategory(ctx echo.Context) error {
	db := kpbatApi.DB()
	bind := models.Category{
		DisplayName: ctx.FormValue("name"),
		Description: ctx.FormValue("description"),
	}
	validated := models.CategoryValidator(&bind)
	if !validated.Ok {
		return view("manage", "Gallery | Categories", ctx, echo.Map{
			"message": validated.Message,
		})
	}
	if err := db.Create(&bind).Error; err != nil {
		return utils.HttpError(ctx, http.StatusInternalServerError, utils.Message(err.Error()))
	}
	if err := utils.CreateDirectory("category_" + strconv.Itoa(bind.ID)); err != nil {
		return utils.HttpError(ctx, http.StatusInternalServerError, utils.Message(err.Error()))
	}
	return view("manage", "Gallery | Categories", ctx, nil)
}
func renderManageCategory(ctx echo.Context) error {
	return view("images", "Gallery | Images", ctx, echo.Map{
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
func setPrimaryImage(ctx echo.Context) error {
	db := kpbatApi.DB()
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		return utils.HttpError(ctx, http.StatusBadRequest, utils.Message("Id param must be integer!"))
	}
	exists, category := services.FindCategory(id)
	if !exists {
		return utils.HttpError(ctx, http.StatusNotFound, utils.Message("Category not found!"))
	}
	category.PrimaryImage = ctx.QueryParam("image")
	db.Save(&category)
	return ctx.Redirect(http.StatusMovedPermanently, fmt.Sprintf("/v1/panel/manage/%d", id))
}
func removeImage(ctx echo.Context) error {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		return utils.HttpError(ctx, http.StatusBadRequest, utils.Message("Id param must be integer!"))
	}
	state := utils.RemoveImage(id, ctx.QueryParam("image"))
	if state {
		services.RemoveImage(ctx.QueryParam("image"))
		return ctx.NoContent(http.StatusOK)
	}
	return ctx.NoContent(http.StatusNotFound)
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
	panel.GET("/manage/:id/remove", removeCategory)
	panel.GET("/manage/:id", renderManageCategory)
	panel.POST("/manage/:id", uploadImages)
	panel.DELETE("/manage/:id", removeImage)
	panel.GET("/manage/:id/primary", setPrimaryImage)
}
