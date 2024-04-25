package database

import (
	"github.com/montheankul-k/assessment-tax/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func DBConnect(c config.IDBConfig) *gorm.DB {
	db, err := gorm.Open(postgres.Open(c.Url()), &gorm.Config{})
	if err != nil {
		panic("Failed to connect to database")
	}

	return db
}
