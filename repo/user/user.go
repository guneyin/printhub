package user

import (
	"context"
	"github.com/guneyin/printhub/market"
	"github.com/guneyin/printhub/model"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
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

func (r *Repo) GetByUUID(ctx context.Context, uuid string) (*model.User, error) {
	ctx = context.WithoutCancel(ctx)

	user := &model.User{}
	tx := r.db.Debug().Where("uuid = ?", uuid).Find(user)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return user, nil
}

func (r *Repo) GetByEmail(ctx context.Context, email string, role model.UserRole) (*model.User, error) {
	ctx = context.WithoutCancel(ctx)

	ur, err := model.NewUserRole(string(role))
	if err != nil {
		return nil, err
	}

	var user *model.User
	tx := r.db.Debug().Model(&model.User{}).Where("email = ? and role = ?", email, ur).First(&user)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return user, nil
}

func (r *Repo) Create(ctx context.Context, u *model.User) error {
	ctx = context.WithoutCancel(ctx)
	tx := r.db.
		Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "role"}, {Name: "email"}},
			UpdateAll: true,
		}).Save(u)

	return tx.Error
}

func (r *Repo) Update(ctx context.Context, u *model.User) error {
	ctx = context.WithoutCancel(ctx)
	tx := r.db.Debug().Updates(u)
	return tx.Error
}

func (r *Repo) migrate() {
	if err := r.db.AutoMigrate(&model.User{}); err != nil {
		panic(err)
	}
}
