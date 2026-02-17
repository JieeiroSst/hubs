package service

import (
	"context"
	"fmt"
	"time"

	"github.com/JIeeiroSst/hub/domain"
	"github.com/JIeeiroSst/hub/repository"
	ws "github.com/JIeeiroSst/hub/websocket"
)


type ItemService struct {
	repo repository.ItemRepository
	hub  *ws.Hub
}

func NewItemService(repo repository.ItemRepository, hub *ws.Hub) *ItemService {
	return &ItemService{repo: repo, hub: hub}
}


func (s *ItemService) Create(ctx context.Context, name, content string) (*domain.Item, error) {
	if name == "" {
		return nil, fmt.Errorf("name is required")
	}

	item := &domain.Item{
		Name:    name,
		Content: content,
	}

	created, err := s.repo.Create(ctx, item)
	if err != nil {
		return nil, fmt.Errorf("create item failed: %w", err)
	}

	s.hub.Broadcast(ws.EventItemCreated, created)

	return created, nil
}

func (s *ItemService) List(ctx context.Context, params domain.ListParams) (*domain.ListResult, error) {
	params.SetDefaults()

	result, err := s.repo.List(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("list items failed: %w", err)
	}
	return result, nil
}

func (s *ItemService) GetByID(ctx context.Context, id string) (*domain.Item, error) {
	item, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("get item failed: %w", err)
	}
	return item, nil
}

func (s *ItemService) HealthCheck(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()
	return s.repo.Ping(ctx)
}
