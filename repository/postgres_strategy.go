package repository

import (
	"context"
	"fmt"
	"math"
	"time"

	"github.com/JIeeiroSst/hub/domain"
	"github.com/google/uuid"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type postgresItem struct {
	ID        string    `gorm:"primaryKey;type:varchar(36)"`
	Name      string    `gorm:"type:varchar(255);not null"`
	Content   string    `gorm:"type:text"`
	CreatedAt time.Time `gorm:"autoCreateTime;index"` 
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
}

func (postgresItem) TableName() string { return "items" }

type PostgresStrategy struct {
	db *gorm.DB
}

func NewPostgresStrategy(dsn string) (*PostgresStrategy, error) {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("postgres connect error: %w", err)
	}

	// Auto migrate tạo table nếu chưa có
	if err := db.AutoMigrate(&postgresItem{}); err != nil {
		return nil, fmt.Errorf("postgres migrate error: %w", err)
	}

	return &PostgresStrategy{db: db}, nil
}

func (r *PostgresStrategy) Create(ctx context.Context, item *domain.Item) (*domain.Item, error) {
	row := &postgresItem{
		ID:      uuid.NewString(),
		Name:    item.Name,
		Content: item.Content,
	}
	if err := r.db.WithContext(ctx).Create(row).Error; err != nil {
		return nil, fmt.Errorf("postgres create error: %w", err)
	}
	return todomainItem(row), nil
}

func (r *PostgresStrategy) List(ctx context.Context, params domain.ListParams) (*domain.ListResult, error) {
	params.SetDefaults()

	var rows []postgresItem
	var total int64

	query := r.db.WithContext(ctx).Model(&postgresItem{})

	if err := query.Count(&total).Error; err != nil {
		return nil, fmt.Errorf("postgres count error: %w", err)
	}

	orderClause := buildOrderClause(params.SortBy, params.SortDir)
	offset := (params.Page - 1) * params.PageSize

	if err := query.
		Order(orderClause).
		Limit(params.PageSize).
		Offset(offset).
		Find(&rows).Error; err != nil {
		return nil, fmt.Errorf("postgres list error: %w", err)
	}

	items := make([]*domain.Item, len(rows))
	for i := range rows {
		items[i] = todomainItem(&rows[i])
	}

	return &domain.ListResult{
		Items:      items,
		Total:      total,
		Page:       params.Page,
		PageSize:   params.PageSize,
		TotalPages: int(math.Ceil(float64(total) / float64(params.PageSize))),
	}, nil
}

func (r *PostgresStrategy) GetByID(ctx context.Context, id string) (*domain.Item, error) {
	var row postgresItem
	if err := r.db.WithContext(ctx).First(&row, "id = ?", id).Error; err != nil {
		return nil, fmt.Errorf("postgres get by id error: %w", err)
	}
	return todomainItem(&row), nil
}

func (r *PostgresStrategy) Ping(ctx context.Context) error {
	sqlDB, err := r.db.DB()
	if err != nil {
		return err
	}
	return sqlDB.PingContext(ctx)
}

func todomainItem(row *postgresItem) *domain.Item {
	return &domain.Item{
		ID:        row.ID,
		Name:      row.Name,
		Content:   row.Content,
		CreatedAt: row.CreatedAt,
		UpdatedAt: row.UpdatedAt,
	}
}
