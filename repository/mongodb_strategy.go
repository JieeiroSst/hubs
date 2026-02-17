package repository

import (
	"context"
	"fmt"
	"math"
	"time"

	"github.com/JIeeiroSst/hub/domain"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoDBStrategy struct {
	collection *mongo.Collection
}

func NewMongoDBStrategy(uri, dbName, collectionName string) (*MongoDBStrategy, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		return nil, fmt.Errorf("mongodb connect error: %w", err)
	}

	if err := client.Ping(ctx, nil); err != nil {
		return nil, fmt.Errorf("mongodb ping error: %w", err)
	}

	collection := client.Database(dbName).Collection(collectionName)

	_, err = collection.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys: bson.D{{Key: "created_at", Value: -1}},
	})
	if err != nil {
		return nil, fmt.Errorf("mongodb create index error: %w", err)
	}

	return &MongoDBStrategy{collection: collection}, nil
}

func (r *MongoDBStrategy) Create(ctx context.Context, item *domain.Item) (*domain.Item, error) {
	now := time.Now()
	doc := &domain.Item{
		ID:        uuid.NewString(),
		Name:      item.Name,
		Content:   item.Content,
		CreatedAt: now,
		UpdatedAt: now,
	}

	if _, err := r.collection.InsertOne(ctx, doc); err != nil {
		return nil, fmt.Errorf("mongodb create error: %w", err)
	}
	return doc, nil
}

func (r *MongoDBStrategy) List(ctx context.Context, params domain.ListParams) (*domain.ListResult, error) {
	params.SetDefaults()

	filter := bson.D{}

	total, err := r.collection.CountDocuments(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("mongodb count error: %w", err)
	}

	sortVal := -1
	if params.SortDir == "asc" {
		sortVal = 1
	}

	skip := int64((params.Page - 1) * params.PageSize)
	limit := int64(params.PageSize)

	opts := options.Find().
		SetSort(bson.D{{Key: params.SortBy, Value: sortVal}}).
		SetSkip(skip).
		SetLimit(limit)

	cursor, err := r.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, fmt.Errorf("mongodb find error: %w", err)
	}
	defer cursor.Close(ctx)

	var items []*domain.Item
	if err := cursor.All(ctx, &items); err != nil {
		return nil, fmt.Errorf("mongodb decode error: %w", err)
	}

	return &domain.ListResult{
		Items:      items,
		Total:      total,
		Page:       params.Page,
		PageSize:   params.PageSize,
		TotalPages: int(math.Ceil(float64(total) / float64(params.PageSize))),
	}, nil
}

func (r *MongoDBStrategy) GetByID(ctx context.Context, id string) (*domain.Item, error) {
	var item domain.Item
	if err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&item); err != nil {
		return nil, fmt.Errorf("mongodb get by id error: %w", err)
	}
	return &item, nil
}

func (r *MongoDBStrategy) Ping(ctx context.Context) error {
	return r.collection.Database().Client().Ping(ctx, nil)
}
