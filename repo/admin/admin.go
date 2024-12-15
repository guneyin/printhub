package admin

import (
	"github.com/guneyin/printhub/market"
	"gorm.io/gorm"
)

type Repo struct {
	db *gorm.DB
}

func NewRepo() *Repo {
	r := &Repo{
		db: market.Get().DB,
	}
	r.migrate()
	return r
}

func (r *Repo) migrate() {
	//panic(r.db.AutoMigrate(&model.User{}))
}
