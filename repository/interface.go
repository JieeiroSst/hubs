package repository

import (
	"context"

	"github.com/JIeeiroSst/hub/domain"
)


type ItemRepository interface {
	Create(ctx context.Context, item *domain.Item) (*domain.Item, error)
	List(ctx context.Context, params domain.ListParams) (*domain.ListResult, error)
	GetByID(ctx context.Context, id string) (*domain.Item, error)
	Ping(ctx context.Context) error
}

type DBType string

const (
	DBTypePostgres DBType = "postgres"
	DBTypeMySQL    DBType = "mysql"
	DBTypeMongoDB  DBType = "mongodb"
)
