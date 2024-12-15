package config

import (
	"context"
	"errors"
	"github.com/guneyin/printhub/market"
	"github.com/guneyin/printhub/model"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Repo struct {
	db *gorm.DB
}

func NewConfigRepo() *Repo {
	r := &Repo{
		db: market.Get().DB,
	}
	r.migrate()
	return r
}

func (r *Repo) Get(ctx context.Context, id, module, key string) (*model.Config, error) {
	ctx = context.WithoutCancel(ctx)
	obj := &model.Config{
		Identifier: id,
		Module:     module,
		Key:        key,
	}

	tx := r.db.Where(obj).First(obj)
	if tx.Error != nil {
		return nil, tx.Error
	}
	if tx.RowsAffected == 0 {
		return nil, errors.New("not found")
	}
	return obj, nil
}

func (r *Repo) Set(ctx context.Context, list *model.ConfigList) error {
	ctx = context.WithoutCancel(ctx)
	tx := r.db.
		Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "identifier"}, {Name: "module"}, {Name: "key"}},
			DoUpdates: clause.AssignmentColumns([]string{"value"}),
		}).Save(list)
	return tx.Error
}

func (r *Repo) Delete(ctx context.Context, id, module, key string) error {
	ctx = context.WithoutCancel(ctx)
	cond := &model.Config{
		Identifier: id,
		Module:     module,
		Key:        key,
	}
	return r.db.Delete(&model.Config{}, cond).Error
}

func (r *Repo) migrate() {
	if err := r.db.AutoMigrate(&model.Config{}); err != nil {
		panic(err)
	}
}
