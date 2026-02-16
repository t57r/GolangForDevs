package handlers

import "context"

type Repository interface {
	PutDocument(ctx context.Context, collection string, doc map[string]any) error
	GetDocument(ctx context.Context, collection string, filter map[string]any) (map[string]any, error)
	ListDocuments(ctx context.Context, collection string, filter map[string]any, limit, skip int64) ([]map[string]any, error)
	DeleteDocument(ctx context.Context, collection string, filter map[string]any) (int64, error)

	CreateCollection(ctx context.Context, name string) error
	ListCollections(ctx context.Context) ([]string, error)
	DeleteCollection(ctx context.Context, name string) error

	CreateIndex(ctx context.Context, collection string, keys map[string]int, unique bool, name string) (string, error)
	DeleteIndex(ctx context.Context, collection string, name string) error
}
