package repository

import (
	"context"
	"fmt"
	"math"
	"time"

	"github.com/JIeeiroSst/hub/domain"
	"github.com/google/uuid"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type mysqlItem struct {
	ID        string    `gorm:"primaryKey;type:varchar(36)"`
	Name      string    `gorm:"type:varchar(255);not null"`
	Content   string    `gorm:"type:text"`
	CreatedAt time.Time `gorm:"autoCreateTime;index"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
}

func (mysqlItem) TableName() string { return "items" }

type MySQLStrategy struct {
	db *gorm.DB
}

func NewMySQLStrategy(dsn string) (*MySQLStrategy, error) {
	// DSN format: "user:pass@tcp(host:3306)/dbname?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("mysql connect error: %w", err)
	}

	if err := db.AutoMigrate(&mysqlItem{}); err != nil {
		return nil, fmt.Errorf("mysql migrate error: %w", err)
	}

	return &MySQLStrategy{db: db}, nil
}

func (r *MySQLStrategy) Create(ctx context.Context, item *domain.Item) (*domain.Item, error) {
	row := &mysqlItem{
		ID:      uuid.NewString(),
		Name:    item.Name,
		Content: item.Content,
	}
	if err := r.db.WithContext(ctx).Create(row).Error; err != nil {
		return nil, fmt.Errorf("mysql create error: %w", err)
	}
	return toMySQLDomain(row), nil
}

func (r *MySQLStrategy) List(ctx context.Context, params domain.ListParams) (*domain.ListResult, error) {
	params.SetDefaults()

	var rows []mysqlItem
	var total int64

	query := r.db.WithContext(ctx).Model(&mysqlItem{})

	if err := query.Count(&total).Error; err != nil {
		return nil, fmt.Errorf("mysql count error: %w", err)
	}

	orderClause := buildOrderClause(params.SortBy, params.SortDir)
	offset := (params.Page - 1) * params.PageSize

	if err := query.
		Order(orderClause).
		Limit(params.PageSize).
		Offset(offset).
		Find(&rows).Error; err != nil {
		return nil, fmt.Errorf("mysql list error: %w", err)
	}

	items := make([]*domain.Item, len(rows))
	for i := range rows {
		items[i] = toMySQLDomain(&rows[i])
	}

	return &domain.ListResult{
		Items:      items,
		Total:      total,
		Page:       params.Page,
		PageSize:   params.PageSize,
		TotalPages: int(math.Ceil(float64(total) / float64(params.PageSize))),
	}, nil
}

func (r *MySQLStrategy) GetByID(ctx context.Context, id string) (*domain.Item, error) {
	var row mysqlItem
	if err := r.db.WithContext(ctx).First(&row, "id = ?", id).Error; err != nil {
		return nil, fmt.Errorf("mysql get by id error: %w", err)
	}
	return toMySQLDomain(&row), nil
}

func (r *MySQLStrategy) Ping(ctx context.Context) error {
	sqlDB, err := r.db.DB()
	if err != nil {
		return err
	}
	return sqlDB.PingContext(ctx)
}

func toMySQLDomain(row *mysqlItem) *domain.Item {
	return &domain.Item{
		ID:        row.ID,
		Name:      row.Name,
		Content:   row.Content,
		CreatedAt: row.CreatedAt,
		UpdatedAt: row.UpdatedAt,
	}
}
