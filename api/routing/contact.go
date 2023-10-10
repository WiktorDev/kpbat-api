package routing

import (
	"github.com/labstack/echo/v4"
	kpbatApi "kpbatApi/api/base"
	"kpbatApi/api/base/utils"
	"net/http"
)

type ContactForm struct {
	Email       string `json:"email"`
	PhoneNumber string `json:"phoneNumber"`
	Subject     string `json:"subject"`
	Message     string `json:"message"`
}

func contactFormValidator(bind *ContactForm) utils.Validated {
	if len(bind.Email) < 6 {
		return utils.Validated{Message: "email length must be longer than 6 characters"}
	}
	if len(bind.Subject) < 12 {
		return utils.Validated{Message: "subject length must be longer than 12 characters"}
	}
	if len(bind.Message) < 20 {
		return utils.Validated{Message: "message length must be longer than 12 characters"}
	}
	return utils.Validated{Ok: true}
}

func sendMail(ctx echo.Context) error {
	var config = kpbatApi.GetConfig()
	bind := new(ContactForm)
	if err := ctx.Bind(bind); err != nil {
		return utils.HttpError(ctx, http.StatusBadRequest, utils.Message(err.Error()))
	}
	err, isValid := utils.Validate(ctx, contactFormValidator(bind))
	if err != nil {
		return err
	}
	if isValid {
		if err := kpbatApi.SendMail(config.Mail.AdminMail, "New message from "+bind.Email, "template.html", bind); err != nil {
			return utils.HttpError(ctx, http.StatusInternalServerError, utils.Message(err.Error()))
		}
	}
	return ctx.NoContent(200)
}

func InitContactRouting(app *echo.Echo) {
	app.POST("/contact", sendMail)
}
