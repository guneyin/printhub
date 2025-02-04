package database

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/guneyin/printhub/config"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	_ "github.com/joho/godotenv/autoload"
)

var gormConfig = &gorm.Config{Logger: logger.Default.LogMode(logger.Error)}

func NewPostgresDB(cfg *config.Config) (*gorm.DB, error) {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Europe/Istanbul",
		cfg.DBHost, cfg.DBUser, cfg.DBPwd, cfg.DBName, cfg.DBPort)
	return gorm.Open(postgres.Open(dsn), gormConfig)
}

func NewSqliteDB(cfg *config.Config) (*gorm.DB, error) {
	path := filepath.Dir(cfg.DBPath)
	if _, err := os.Stat(path); err != nil {
		err = os.MkdirAll(path, os.ModePerm)
		if err != nil {
			return nil, err
		}
	}
	return gorm.Open(sqlite.Open(cfg.DBPath), gormConfig)
}

func NewTestDB() (*gorm.DB, error) {
	return gorm.Open(sqlite.Open("file::memory:?cache=shared"), gormConfig)
}
