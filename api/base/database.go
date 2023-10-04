package kpbatApi

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	kpbatApi "kpbatApi/api/models"
)

var database *gorm.DB

func InitDatabase(config DatabaseConfig) *gorm.DB {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local", config.Username, config.Password, config.Hostname, config.Port, config.Database)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	database = db
	db.AutoMigrate(&kpbatApi.Category{}, &kpbatApi.Image{})
	return db
}

func DB() *gorm.DB {
	return database
}
