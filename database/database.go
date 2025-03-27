package database

import (
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"manga-catalog/models"
)

var DB *gorm.DB

func ConnectDB() {
	dsn := "host=localhost user=postgres password=2705 dbname=manga port=5432 sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Ошибка подключения к БД:", err)
	}

	DB = db

	err = db.AutoMigrate(&models.Manga{})
	if err != nil {
		log.Fatal("Ошибка миграции:", err)
	}

	fmt.Println("Подключение к БД установлено")
}
