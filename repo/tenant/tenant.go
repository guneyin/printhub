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

func (r *Repo) Save(ctx context.Context, t *model.Tenant) (*model.Tenant, error) {
	ctx = context.WithoutCancel(ctx)
	tx := r.db.WithContext(ctx).Save(t)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return t, nil
}

func (r *Repo) Search(ctx context.Context, filter model.QueryFilter) (model.TenantList, error) {
	ctx = context.WithoutCancel(ctx)
	tenant := model.TenantList{}
	q, args := filter.Query()
	tx := r.db.
		Debug().
		WithContext(ctx).Model(&model.Tenant{}).
		Where(q, args...).
		Find(&tenant)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return tenant, nil
}

func (r *Repo) GetUserList(ctx context.Context, id string) (model.UserList, error) {
	ctx = context.WithoutCancel(ctx)
	userList := model.UserList{}
	tx := r.db.WithContext(ctx).Model(&model.TenantUser{}).
		InnerJoins("users").
		Where("tenant_id = ?", id).Find(&userList)
	if tx.Error != nil {
		return nil, tx.Error
	}

	return userList, nil
}

func (r *Repo) AddUser(ctx context.Context, t *model.Tenant, u *model.User) error {
	ctx = context.WithoutCancel(ctx)
	tx := r.db.WithContext(ctx)

	return tx.Save(&model.TenantUser{
		TenantID: t.ID,
		UserID:   u.ID,
	}).Error
}
