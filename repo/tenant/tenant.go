package tenant

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

func (r *Repo) migrate() {
	if err := r.db.AutoMigrate(&model.Tenant{}, &model.TenantUser{}); err != nil {
		panic(err)
	}
}

func (r *Repo) GetByUUID(ctx context.Context, uuid string) (*model.Tenant, error) {
	ctx = context.WithoutCancel(ctx)
	tx := r.db.WithContext(ctx)

	tenant := &model.Tenant{UUID: uuid}
	err := tx.First(tenant).Error
	return tenant, err
}

func (r *Repo) AddUser(ctx context.Context, t *model.Tenant, u *model.User) error {
	ctx = context.WithoutCancel(ctx)
	tx := r.db.WithContext(ctx)

	return tx.Save(&model.TenantUser{
		TenantID: t.ID,
		UserID:   u.ID,
	}).Error
}
