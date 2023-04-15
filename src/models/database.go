package models

import (
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDatabase() (*gorm.DB, error) {
	dsn := "host=localhost user=taskmanager password=distributed-tasks dbname=taskmanagementdb port=5432 sslmode=disable TimeZone=Brazil/East"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatal(err.Error())
	}

	if err = db.AutoMigrate(&Task{}); err != nil {
		log.Println(err)
	}

	return db, err
}
