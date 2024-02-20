package config

import (
	"fmt"
	"os"

	reportModel "gdsc/baro/app/report/models"
	userModel "gdsc/baro/app/user/models"
	videoModel "gdsc/baro/app/video/models"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDatabase() (*gorm.DB, error) {
	MYSQL_HOST := os.Getenv("MYSQL_HOST")
	MYSQL_PORT := os.Getenv("MYSQL_PORT")
	MYSQL_USER := os.Getenv("MYSQL_USER")
	MYSQL_PASSWORD := os.Getenv("MYSQL_PASSWORD")
	MYSQL_DATABASE := os.Getenv("MYSQL_DATABASE")

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Asia%%2FSeoul",
		MYSQL_USER, MYSQL_PASSWORD, MYSQL_HOST, MYSQL_PORT, MYSQL_DATABASE)

	database, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		panic("Failed to connect to database!")
	}

	userErr := database.AutoMigrate(&userModel.User{})
	if userErr != nil {
		return nil, userErr
	}

	reportErr := database.AutoMigrate(&reportModel.Report{})
	if reportErr != nil {
		return nil, reportErr
	}

	videoErr := database.AutoMigrate(&videoModel.Video{})
	if videoErr != nil {
		return nil, videoErr
	}

	DB = database

	return DB, nil
}
