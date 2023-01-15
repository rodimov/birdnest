package app

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func Connect(dbConfig *Config) (*gorm.DB, error) {
	return gorm.Open(postgres.Open(dbConfig.Datasource), &gorm.Config{})
}
