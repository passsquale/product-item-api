package repo

import (
	"context"
	"github.com/jmoiron/sqlx"
	"github.com/passsquale/product-item-api/internal/model"
)

// Repo is DAO for Template
type Repo interface {
	DescribeTemplate(ctx context.Context, templateID uint64) (*model.Item, error)
}

type repo struct {
	db        *sqlx.DB
	batchSize uint
}

// NewRepo returns Repo interface
func NewRepo(db *sqlx.DB, batchSize uint) Repo {
	return &repo{db: db, batchSize: batchSize}
}

func (r *repo) DescribeTemplate(ctx context.Context, itemID uint64) (*model.Item, error) {
	return nil, nil
}
