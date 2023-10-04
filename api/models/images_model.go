package models

type Image struct {
	ID         int    `gorm:"primaryKey" json:"id"`
	CategoryID int    `json:"category_id"`
	FileName   string `gorm:"size:64" json:"file_name"`
}

func (Image) TableName() string {
	return "images"
}
