package models

import (
	"kpbatApi/api/base/utils"
	"regexp"
)

type ContactForm struct {
	Email       string `json:"email"`
	PhoneNumber string `json:"phone_number"`
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	Subject     string `json:"subject"`
	Message     string `json:"message"`
}

func ContactFormValidator(bind *ContactForm) utils.Validated {
	if !isValidPattern("^[\\w-\\.]+@([\\w-]+\\.)+[\\w-]{2,4}$", bind.Email) {
		return utils.Validated{Message: "Invalid email format!"}
	}
	if len(bind.FirstName) < 3 {
		return utils.Validated{Message: "Invalid first name!"}
	}
	if len(bind.LastName) < 3 {
		return utils.Validated{Message: "Invalid last name!"}
	}
	if len(bind.Message) < 20 {
		return utils.Validated{Message: "message length must be longer than 12 characters"}
	}
	return utils.Validated{Ok: true}
}

func isValidPattern(pattern string, text string) bool {
	match, _ := regexp.MatchString(pattern, text)
	return match
}
