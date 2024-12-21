package admin

import (
	"context"
	"github.com/guneyin/printhub/market"
	"github.com/guneyin/printhub/model"
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

func (r *Repo) GetCount(ctx context.Context) (int64, error) {
	var cnt int64
	tx := r.db.Debug().WithContext(ctx).Model(&model.User{}).Where("role = ?", model.UserRoleAdmin)
	tx.Count(&cnt)
	return cnt, tx.Error
}

func (r *Repo) Boostrap(ctx context.Context, u *model.User) error {
	tx := r.db.WithContext(ctx).Create(u)
	return tx.Error
}

func (r *Repo) migrate() {}
