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
		db: market.Get().DB.Debug(),
	}
	r.migrate()
	return r
}

func (r *Repo) GetCount(ctx context.Context) (int64, error) {
	var cnt int64
	tx := r.db.WithContext(ctx).Model(&model.User{}).Where("role = ?", model.UserRoleAdmin)
	tx.Count(&cnt)
	return cnt, tx.Error
}

func (r *Repo) GetTenantList(ctx context.Context) (model.TenantList, error) {
	ctx = context.WithoutCancel(ctx)
	list := model.TenantList{}
	tx := r.db.WithContext(ctx).Model(&model.TenantList{}).Find(&list)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return list, nil
}

func (r *Repo) GetTenantById(ctx context.Context, id string) (*model.Tenant, error) {
	ctx = context.WithoutCancel(ctx)
	tenant := &model.Tenant{}
	tx := r.db.WithContext(ctx).Model(&model.Tenant{}).Where("uuid = ?", id).First(&tenant)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return tenant, nil
}

func (r *Repo) CreateTenant(ctx context.Context, t *model.Tenant) (*model.Tenant, error) {
	ctx = context.WithoutCancel(ctx)
	tx := r.db.WithContext(ctx).Create(t)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return t, nil
}

func (r *Repo) Boostrap(ctx context.Context, u *model.User) error {
	tx := r.db.WithContext(ctx).Create(u)
	return tx.Error
}

func (r *Repo) migrate() {}
