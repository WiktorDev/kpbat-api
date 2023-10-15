package dto

import "kpbatApi/api/base/utils"

type ContactForm struct {
	Email       string `json:"email"`
	PhoneNumber string `json:"phoneNumber"`
	Subject     string `json:"subject"`
	Message     string `json:"message"`
}

func ContactFormValidator(bind *ContactForm) utils.Validated {
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
