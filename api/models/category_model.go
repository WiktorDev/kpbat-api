package models

import "kpbatApi/api/base/utils"

type Category struct {
	ID           int     `gorm:"primaryKey" json:"id"`
	DisplayName  string  `json:"display_name"`
	Description  string  `json:"description"`
	PrimaryImage string  `json:"primary_image"`
	Images       []Image `json:"images"`
}

func (Category) TableName() string {
	return "categories"
}

func CategoryValidator(bind *Category) utils.Validated {
	if len(bind.DisplayName) < 6 {
		return utils.Validated{Message: "Category name length must be longer than 6 characters"}
	}
	return utils.Validated{Ok: true}
}
