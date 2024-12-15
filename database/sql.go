package database

import (
	"fmt"
	"github.com/guneyin/printhub/config"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"os"
	"path/filepath"

	_ "github.com/joho/godotenv/autoload"
)

var gormConfig = &gorm.Config{Logger: logger.Default.LogMode(logger.Error)}

func NewPostgresDB(cfg *config.Config) (*gorm.DB, error) {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Europe/Istanbul",
		cfg.DbHost, cfg.DbUser, cfg.DbPwd, cfg.DbName, cfg.DbPort)
	return gorm.Open(postgres.Open(dsn), gormConfig)
}

func NewSqliteDB(cfg *config.Config) (*gorm.DB, error) {
	path := filepath.Dir(cfg.DbPath)
	if _, err := os.Stat(path); err != nil {
		err = os.MkdirAll(path, os.ModePerm)
		if err != nil {
			return nil, err
		}
	}
	return gorm.Open(sqlite.Open(cfg.DbPath), gormConfig)
}

func NewTestDB() (*gorm.DB, error) {
	return gorm.Open(sqlite.Open("file::memory:?cache=shared"), gormConfig)
}
