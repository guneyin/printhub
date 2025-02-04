package market

import (
	"log/slog"
	"os"
	"sync"

	"github.com/guneyin/printhub/config"
	"github.com/guneyin/printhub/database"
	"gorm.io/gorm"
)

var (
	once   sync.Once
	market *Market
)

type Market struct {
	Config *config.Config
	DB     *gorm.DB
	Log    *slog.Logger
}

func InitMarket() {
	once.Do(func() {
		cfg, err := config.New()
		handleErr(err)

		db, err := database.NewSqliteDB(cfg)
		handleErr(err)

		logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

		market = &Market{
			Config: cfg,
			DB:     db,
			Log:    logger,
		}
	})
}

func InitTestMarket() {
	once.Do(func() {
		cfg, err := config.New()
		handleErr(err)

		db, err := database.NewSqliteDB(cfg)
		handleErr(err)

		logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

		market = &Market{
			Config: cfg,
			DB:     db,
			Log:    logger,
		}
	})
}

func Get() *Market {
	return market
}

func Log() *slog.Logger {
	return market.Log
}

func handleErr(e error) {
	if e != nil {
		panic(e)
	}
}
