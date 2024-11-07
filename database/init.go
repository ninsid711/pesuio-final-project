package database

import (
	"github.com/anuragrao04/pesuio-final-project/models"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Init(databaseFileName string) {
	// implement
	// populate DB variable
	var err error
	DB, err := gorm.Open(sqlite.Open(databaseFileName), &gorm.Config{})
	if err != nil {
		panic("db connection failed")
	}
	DB.AutoMigrate(&models.User, &models.Question, &models.TestCase)
	return DB
}
