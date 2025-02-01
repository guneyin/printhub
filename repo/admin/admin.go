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
	tx := r.db.WithContext(ctx).Model(&model.User{}).Where("role = ?", model.UserRoleAdmin)
	tx.Count(&cnt)
	return cnt, tx.Error
}

//func (r *Repo) GetTenantList(ctx context.Context) (model.TenantList, error) {
//	ctx = context.WithoutCancel(ctx)
//	list := model.TenantList{}
//	tx := r.db.WithContext(ctx).Model(&model.TenantList{}).Find(&list)
//	if tx.Error != nil {
//		return nil, tx.Error
//	}
//	return list, nil
//}

//func (r *Repo) GetTenant(ctx context.Context, filter model.QueryFilter) (model.TenantList, error) {
//	ctx = context.WithoutCancel(ctx)
//	tenant := model.TenantList{}
//	q, args := filter.Query()
//	tx := r.db.
//		Debug().
//		WithContext(ctx).Model(&model.Tenant{}).
//		Where(q, args...).
//		Find(&tenant)
//	if tx.Error != nil {
//		return nil, tx.Error
//	}
//	return tenant, nil
//}

//func (r *Repo) SaveTenant(ctx context.Context, t *model.Tenant) (*model.Tenant, error) {
//	ctx = context.WithoutCancel(ctx)
//	tx := r.db.WithContext(ctx).Save(t)
//	if tx.Error != nil {
//		return nil, tx.Error
//	}
//	return t, nil
//}

//func (r *Repo) GetTenantUsers(ctx context.Context, id string) ([]model.User, error) {
//	ctx = context.WithoutCancel(ctx)
//	var users []model.User
//
//}

func (r *Repo) Boostrap(ctx context.Context, u *model.User) error {
	tx := r.db.WithContext(ctx).Create(u)
	return tx.Error
}

func (r *Repo) migrate() {}
