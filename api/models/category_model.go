package models

import http_utils "kpbatApi/api/base/utils"

type Category struct {
	ID              int    `gorm:"primaryKey" json:"id"`
	DisplayName     string `json:"display_name"`
	Description     string `json:"description"`
	ImagesDirectory string `json:"images_directory"`
}

func (Category) TableName() string {
	return "categories"
}

func CategoryValidator(bind *Category) http_utils.Validated {
	if len(bind.DisplayName) < 6 {
		return http_utils.Validated{Message: "display name length must be longer than 6 characters"}
	}
	if len(bind.Description) < 12 {
		return http_utils.Validated{Message: "description length must be longer than 12 characters"}
	}
	if len(bind.ImagesDirectory) < 2 {
		return http_utils.Validated{Message: "images directory length must be longer than 6 characters"}
	}
	return http_utils.Validated{Ok: true}
}
