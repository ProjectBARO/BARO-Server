package models

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDatabase() {
	DB_USER := os.Getenv("DB_USER")
	DB_PASS := os.Getenv("DB_PASSWORD")
	DB_NAME := os.Getenv("DB_NAME")
	INSTANCE_CONNECTION_NAME := os.Getenv("INSTANCE_CONNECTION_NAME")

	dbURI := fmt.Sprintf("%s:%s@unix(/cloudsql/%s)/%s?parseTime=true",
		DB_USER, DB_PASS, INSTANCE_CONNECTION_NAME, DB_NAME)

	dbPool, err := sql.Open("mysql", dbURI)
	if err != nil {
		log.Fatalf("Could not connect to the database: %v", err)
	}

	database, err := gorm.Open(mysql.New(mysql.Config{
		Conn: dbPool,
	}), &gorm.Config{})

	if err != nil {
		log.Fatalf("Could not connect to the database: %v", err)
	}

	DB = database
}
