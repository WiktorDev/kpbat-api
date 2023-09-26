package kpbatApi

type Category struct {
	ID              int    `gorm:"primaryKey" json:"id"`
	DisplayName     string `json:"display_name"`
	Description     string `json:"description"`
	ImagesDirectory string `json:"images_directory"`
}

func (Category) TableName() string {
	return "categories"
}
