package market

import (
	"github.com/guneyin/printhub/config"
	"github.com/guneyin/printhub/database"
	"gorm.io/gorm"
	"sync"
)

var (
	once   sync.Once
	market *Market
)

type Market struct {
	Config *config.Config
	DB     *gorm.DB
}

func InitMarket() {
	once.Do(func() {
		cfg, err := config.New()
		handleErr(err)

		db, err := database.NewSqliteDB(cfg)
		handleErr(err)

		market = &Market{
			Config: cfg,
			DB:     db,
		}
	})
}

func InitTestMarket() {
	once.Do(func() {
		cfg, err := config.New()
		handleErr(err)

		db, err := database.NewSqliteDB(cfg)
		handleErr(err)

		market = &Market{
			Config: cfg,
			DB:     db,
		}
	})
}

func Get() *Market {
	return market
}

func handleErr(e error) {
	if e != nil {
		panic(e)
	}
}
