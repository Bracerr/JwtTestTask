package database

import (
	"JwtTestTask/src/pkg/logger"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type DbInitModel struct {
	DbHost     string
	DbUser     string
	DbPassword string
	DbName     string
	DbPort     string
}

func NewClient(dbModel DbInitModel) *gorm.DB {
	dsn := "host=" + dbModel.DbHost +
		" user=" + dbModel.DbUser +
		" password=" + dbModel.DbPassword +
		" dbname=" + dbModel.DbName +
		" port=" + dbModel.DbPort +
		" sslmode=disable"

	db, dbErr := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if dbErr != nil {
		logger.Log.Fatal("Failed to connect to database")
	}

	return db
}
