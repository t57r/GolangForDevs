package handlers

import "context"

// Used for mocking in unit tests
type FakeRepository struct {
	PutDocumentFn      func(ctx context.Context, collection string, doc map[string]any) error
	GetDocumentFn      func(ctx context.Context, collection string, filter map[string]any) (map[string]any, error)
	ListDocumentsFn    func(ctx context.Context, collection string, filter map[string]any, limit, skip int64) ([]map[string]any, error)
	DeleteDocumentFn   func(ctx context.Context, collection string, filter map[string]any) (int64, error)
	CreateCollectionFn func(ctx context.Context, name string) error
	ListCollectionsFn  func(ctx context.Context) ([]string, error)
	DeleteCollectionFn func(ctx context.Context, name string) error
	CreateIndexFn      func(ctx context.Context, collection string, keys map[string]int, unique bool, name string) (string, error)
	DeleteIndexFn      func(ctx context.Context, collection string, name string) error
}

func (f *FakeRepository) PutDocument(ctx context.Context, collection string, doc map[string]any) error {
	return f.PutDocumentFn(ctx, collection, doc)
}

func (f *FakeRepository) GetDocument(ctx context.Context, collection string, filter map[string]any) (map[string]any, error) {
	return f.GetDocumentFn(ctx, collection, filter)
}

func (f *FakeRepository) ListDocuments(ctx context.Context, collection string, filter map[string]any, limit, skip int64) ([]map[string]any, error) {
	return f.ListDocumentsFn(ctx, collection, filter, limit, skip)
}

func (f *FakeRepository) DeleteDocument(ctx context.Context, collection string, filter map[string]any) (int64, error) {
	return f.DeleteDocumentFn(ctx, collection, filter)
}

func (f *FakeRepository) CreateCollection(ctx context.Context, name string) error {
	return f.CreateCollectionFn(ctx, name)
}

func (f *FakeRepository) ListCollections(ctx context.Context) ([]string, error) {
	return f.ListCollectionsFn(ctx)
}

func (f *FakeRepository) DeleteCollection(ctx context.Context, name string) error {
	return f.DeleteCollectionFn(ctx, name)
}

func (f *FakeRepository) CreateIndex(ctx context.Context, collection string, keys map[string]int, unique bool, name string) (string, error) {
	return f.CreateIndexFn(ctx, collection, keys, unique, name)
}

func (f *FakeRepository) DeleteIndex(ctx context.Context, collection string, name string) error {
	return f.DeleteIndexFn(ctx, collection, name)
}
