package repository

import (
	"fmt"
)

type DBConfig struct {
	Type          DBType
	DSN           string // dùng cho Postgres, MySQL
	MongoURI      string // dùng cho MongoDB
	MongoDBName   string
	MongoCollName string
}


func NewRepository(cfg DBConfig) (ItemRepository, error) {
	switch cfg.Type {
	case DBTypePostgres:
		return NewPostgresStrategy(cfg.DSN)

	case DBTypeMySQL:
		return NewMySQLStrategy(cfg.DSN)

	case DBTypeMongoDB:
		return NewMongoDBStrategy(cfg.MongoURI, cfg.MongoDBName, cfg.MongoCollName)

	default:
		return nil, fmt.Errorf("unsupported db type: %s", cfg.Type)
	}
}

func buildOrderClause(sortBy, sortDir string) string {
	allowedFields := map[string]string{
		"created_at": "created_at",
		"updated_at": "updated_at",
		"name":       "name",
		"id":         "id",
	}

	field, ok := allowedFields[sortBy]
	if !ok {
		field = "created_at" 
	}

	dir := "DESC"
	if sortDir == "asc" {
		dir = "ASC"
	}

	return fmt.Sprintf("%s %s", field, dir)
}
