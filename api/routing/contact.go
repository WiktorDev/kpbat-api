package routing

import (
	"github.com/labstack/echo/v4"
	kpbatApi "kpbatApi/api/base"
	"kpbatApi/api/base/utils"
	"kpbatApi/api/models"
	"net/http"
)

// @Summary	Send message to kpbat.com management
// @Tags		contact
// @Accept		json
// @Produce	json
// @Param		request	body	models.ContactForm	true	"body"
// @Success	200
// @Failure	400	{object}	utils.MessageStruct
// @Failure	500	{object}	utils.MessageStruct
// @Router		/contact [post]
func sendMail(ctx echo.Context) error {
	var config = kpbatApi.GetConfig()
	bind := new(models.ContactForm)
	if err := (&echo.DefaultBinder{}).BindBody(ctx, bind); err != nil {
		return utils.HttpError(ctx, http.StatusInternalServerError, utils.Message(err.Error()))
	}
	err, isValid := utils.Validate(ctx, models.ContactFormValidator(bind))
	if err != nil {
		return err
	}
	if isValid {
		if err := kpbatApi.SendMail(bind.Email, "Message sent!", "template_to_sender.html", nil); err != nil {
			return utils.HttpError(ctx, http.StatusInternalServerError, utils.Message(err.Error()))
		}
		if err := kpbatApi.SendMail(config.Mail.AdminMail, "New message from "+bind.Email, "template_to_system.html", bind); err != nil {
			return utils.HttpError(ctx, http.StatusInternalServerError, utils.Message(err.Error()))
		}
	}
	return ctx.NoContent(200)
}

func InitContactRouting(app *echo.Echo) {
	app.POST("/contact", sendMail)
}
