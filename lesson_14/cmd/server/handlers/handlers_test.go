package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
)

// PUT DOCUMENT -------------------------------------------------------

func TestHandlePutDocument_OK(t *testing.T) {
	requiredCollection := "users"
	requiredUserID := "123"
	requiredName := "Alex"

	repo := &FakeRepository{
		PutDocumentFn: func(ctx context.Context, collection string, doc map[string]any) error {
			if collection != requiredCollection {
				t.Fatalf("wrong collection: %s", collection)
			}
			if doc["user_id"] != requiredUserID {
				t.Fatalf("wrong user_id: %v", doc["user_id"])
			}
			if doc["name"] != requiredName {
				t.Fatalf("wrong name: %v", doc["name"])
			}
			return nil
		},
	}

	srv := NewServer(repo)
	app := fiber.New()
	app.Post("/put_document", srv.HandlePutDocument)

	payload := PutDocumentRequest{
		CollectionName: requiredCollection,
		Document: map[string]any{
			"user_id": requiredUserID,
			"name":    requiredName,
		},
	}
	req := httptest.NewRequest("POST", "/put_document", bytes.NewReader(mustJSON(t, payload)))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("app.Test: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		t.Fatalf("status=%d", resp.StatusCode)
	}

	var got map[string]any
	_ = json.NewDecoder(resp.Body).Decode(&got)
	if got["ok"] != true {
		t.Fatalf("ok=%v", got["ok"])
	}
}

func TestHandlePutDocument_RepoError(t *testing.T) {
	repo := &FakeRepository{
		PutDocumentFn: func(ctx context.Context, collection string, doc map[string]any) error {
			return errors.New("boom")
		},
	}
	srv := NewServer(repo)

	app := fiber.New()
	app.Post("/put_document", srv.HandlePutDocument)

	payload := PutDocumentRequest{
		CollectionName: "users",
		Document:       map[string]any{"user_id": "1"},
	}
	req := httptest.NewRequest("POST", "/put_document", bytes.NewReader(mustJSON(t, payload)))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("app.Test: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != fiber.StatusInternalServerError {
		t.Fatalf("status=%d want 500", resp.StatusCode)
	}
}

// GET DOCUMENT -------------------------------------------------------

func TestHandleGetDocument_OK(t *testing.T) {
	requiredCollection := "users"
	requiredUserID := "123"
	requiredUserName := "Alex"

	repo := &FakeRepository{
		GetDocumentFn: func(ctx context.Context, collection string, filter map[string]any) (map[string]any, error) {
			if collection != requiredCollection {
				t.Fatalf("wrong collection")
			}
			if filter["user_id"] != requiredUserID {
				t.Fatalf("wrong filter")
			}
			return map[string]any{
				"ok":    true,
				"found": true,
				"document": map[string]any{
					"user_id": requiredUserID,
					"name":    requiredUserName,
				},
			}, nil
		},
	}

	srv := NewServer(repo)

	app := fiber.New()
	app.Post("/get_document", srv.HandleGetDocument)

	payload := GetDocumentRequest{
		CollectionName: requiredCollection,
		Filter: map[string]any{
			"user_id": requiredUserID,
		},
	}
	req := httptest.NewRequest("POST", "/get_document", bytes.NewReader(mustJSON(t, payload)))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("app.Test: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		t.Fatalf("status = %d", resp.StatusCode)
	}

	var got map[string]any
	_ = json.NewDecoder(resp.Body).Decode(&got)

	if got["ok"] != true {
		t.Fatalf("ok=%v", got["ok"])
	}
	if got["found"] != true {
		t.Fatalf("found=%v", got["found"])
	}
	foundDoc := got["document"].(map[string]any)
	if foundDoc == nil {
		t.Fatalf("document not found")
	}
	if foundDoc["user_id"] != requiredUserID {
		t.Fatalf("foundDoc[user_id]=%v, required(%s)", got["user_id"], requiredUserID)
	}
	if foundDoc["name"] != requiredUserName {
		t.Fatalf("foundDoc[name]=%v, required(%s)", got["name"], requiredUserName)
	}
}

// LIST DOCUMENTS -----------------------------------------------------

func TestHandleListDocuments_OK(t *testing.T) {
	requiredCollection := "users"
	requiredFilter := map[string]any{"active": true}
	var requiredLimit int64 = 10
	var requiredSkip int64 = 5

	repo := &FakeRepository{
		ListDocumentsFn: func(ctx context.Context, collection string, filter map[string]any, limit, skip int64) ([]map[string]any, error) {
			if collection != requiredCollection {
				t.Fatalf("wrong collection")
			}
			if filter["active"] != true {
				t.Fatalf("wrong filter")
			}
			if limit != requiredLimit {
				t.Fatalf("wrong limit: %d", limit)
			}
			if skip != requiredSkip {
				t.Fatalf("wrong skip: %d", skip)
			}
			return []map[string]any{
				{"user_id": "1", "name": "A"},
				{"user_id": "2", "name": "B"},
			}, nil
		},
	}

	srv := NewServer(repo)
	app := fiber.New()
	app.Post("/list_documents", srv.HandleListDocuments)

	payload := ListDocumentsRequest{
		CollectionName: requiredCollection,
		Filter:         requiredFilter,
		Limit:          requiredLimit,
		Skip:           requiredSkip,
	}
	req := httptest.NewRequest("POST", "/list_documents", bytes.NewReader(mustJSON(t, payload)))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("app.Test: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		t.Fatalf("status=%d", resp.StatusCode)
	}

	var got map[string]any
	_ = json.NewDecoder(resp.Body).Decode(&got)
	if got["ok"] != true {
		t.Fatalf("ok=%v", got["ok"])
	}
	docs, ok := got["documents"].([]any)
	if !ok {
		t.Fatalf("documents type=%T", got["documents"])
	}
	if len(docs) != 2 {
		t.Fatalf("documents len=%d", len(docs))
	}
}

// DELETE DOCUMENT ----------------------------------------------------

func TestHandleDeleteDocument_OK(t *testing.T) {
	requiredCollection := "users"
	requiredUserID := "123"
	var requiredDeleted int64 = 1

	repo := &FakeRepository{
		DeleteDocumentFn: func(ctx context.Context, collection string, filter map[string]any) (int64, error) {
			if collection != requiredCollection {
				t.Fatalf("wrong collection")
			}
			if filter["user_id"] != requiredUserID {
				t.Fatalf("wrong filter user_id")
			}
			return requiredDeleted, nil
		},
	}

	srv := NewServer(repo)
	app := fiber.New()
	app.Post("/delete_document", srv.HandleDeleteDocument)

	payload := DeleteDocumentRequest{
		CollectionName: requiredCollection,
		Filter:         map[string]any{"user_id": requiredUserID},
	}
	req := httptest.NewRequest("POST", "/delete_document", bytes.NewReader(mustJSON(t, payload)))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("app.Test: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		t.Fatalf("status=%d", resp.StatusCode)
	}

	var got map[string]any
	_ = json.NewDecoder(resp.Body).Decode(&got)
	if got["ok"] != true {
		t.Fatalf("ok=%v", got["ok"])
	}

	// JSON decodes numbers as float64 when decoding into map[string]any
	if got["deleted"].(float64) != float64(requiredDeleted) {
		t.Fatalf("deleted=%v", got["deleted"])
	}
}

// CREATE COLLECTION --------------------------------------------------

func TestHandleCreateCollection_OK(t *testing.T) {
	requiredName := "users"

	repo := &FakeRepository{
		CreateCollectionFn: func(ctx context.Context, name string) error {
			if name != requiredName {
				t.Fatalf("wrong name: %s", name)
			}
			return nil
		},
	}

	srv := NewServer(repo)
	app := fiber.New()
	app.Post("/create_collection", srv.HandleCreateCollection)

	payload := CreateCollectionRequest{CollectionName: requiredName}
	req := httptest.NewRequest("POST", "/create_collection", bytes.NewReader(mustJSON(t, payload)))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("app.Test: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		t.Fatalf("status=%d", resp.StatusCode)
	}

	var got map[string]any
	_ = json.NewDecoder(resp.Body).Decode(&got)
	if got["ok"] != true {
		t.Fatalf("ok=%v", got["ok"])
	}
}

// LIST COLLECTIONS ---------------------------------------------------

func TestHandleListCollections_OK(t *testing.T) {
	repo := &FakeRepository{
		ListCollectionsFn: func(ctx context.Context) ([]string, error) {
			return []string{"users", "products"}, nil
		},
	}

	srv := NewServer(repo)
	app := fiber.New()
	app.Post("/list_collections", srv.HandleListCollections)

	req := httptest.NewRequest("POST", "/list_collections", bytes.NewReader([]byte(`{}`)))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("app.Test: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		t.Fatalf("status=%d", resp.StatusCode)
	}

	var got map[string]any
	_ = json.NewDecoder(resp.Body).Decode(&got)

	if got["ok"] != true {
		t.Fatalf("ok=%v", got["ok"])
	}
	colls, ok := got["collections"].([]any)
	if !ok || len(colls) != 2 {
		t.Fatalf("collections=%T len=%d", got["collections"], len(colls))
	}
}

// DELETE COLLECTION --------------------------------------------------

func TestHandleDeleteCollection_OK(t *testing.T) {
	required := "users"
	repo := &FakeRepository{
		DeleteCollectionFn: func(ctx context.Context, name string) error {
			if name != required {
				t.Fatalf("wrong name: %s", name)
			}
			return nil
		},
	}

	srv := NewServer(repo)
	app := fiber.New()
	app.Post("/delete_collection", srv.HandleDeleteCollection)

	payload := DeleteCollectionRequest{CollectionName: required}
	req := httptest.NewRequest("POST", "/delete_collection", bytes.NewReader(mustJSON(t, payload)))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("app.Test: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		t.Fatalf("status=%d", resp.StatusCode)
	}

	var got map[string]any
	_ = json.NewDecoder(resp.Body).Decode(&got)
	if got["ok"] != true {
		t.Fatalf("ok=%v", got["ok"])
	}
}

// CREATE INDEX -------------------------------------------------------

func TestHandleCreateIndex_OK(t *testing.T) {
	requiredCollection := "users"

	repo := &FakeRepository{
		CreateIndexFn: func(ctx context.Context, collection string, keys map[string]int, unique bool, name string) (string, error) {
			if collection != requiredCollection {
				t.Fatalf("wrong collection")
			}
			if keys["user_id"] != 1 {
				t.Fatalf("wrong keys: %v", keys)
			}
			if unique != true {
				t.Fatalf("wrong unique: %v", unique)
			}
			if name != "uidx" {
				t.Fatalf("wrong name: %s", name)
			}
			return "uidx", nil
		},
	}

	srv := NewServer(repo)
	app := fiber.New()
	app.Post("/create_index", srv.HandleCreateIndex)

	payload := CreateIndexRequest{
		CollectionName: requiredCollection,
		Keys:           map[string]int{"user_id": 1},
		Unique:         true,
		Name:           "uidx",
	}
	req := httptest.NewRequest("POST", "/create_index", bytes.NewReader(mustJSON(t, payload)))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("app.Test: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		t.Fatalf("status=%d", resp.StatusCode)
	}

	var got map[string]any
	_ = json.NewDecoder(resp.Body).Decode(&got)
	if got["ok"] != true {
		t.Fatalf("ok=%v", got["ok"])
	}
	if got["index_name"] != "uidx" {
		t.Fatalf("index_name=%v", got["index_name"])
	}
}

// DELETE INDEX -------------------------------------------------------

func TestHandleDeleteIndex_OK(t *testing.T) {
	requiredCollection := "users"
	requiredIndexName := "uidx"

	repo := &FakeRepository{
		DeleteIndexFn: func(ctx context.Context, collection string, name string) error {
			if collection != requiredCollection {
				t.Fatalf("wrong collection")
			}
			if name != requiredIndexName {
				t.Fatalf("wrong index name")
			}
			return nil
		},
	}

	srv := NewServer(repo)
	app := fiber.New()
	app.Post("/delete_index", srv.HandleDeleteIndex)

	payload := DeleteIndexRequest{
		CollectionName: requiredCollection,
		Name:           requiredIndexName,
	}
	req := httptest.NewRequest("POST", "/delete_index", bytes.NewReader(mustJSON(t, payload)))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("app.Test: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		t.Fatalf("status=%d", resp.StatusCode)
	}

	var got map[string]any
	_ = json.NewDecoder(resp.Body).Decode(&got)
	if got["ok"] != true {
		t.Fatalf("ok=%v", got["ok"])
	}
}

// Helpers

func mustJSON(t *testing.T, v any) []byte {
	t.Helper()
	b, err := json.Marshal(v)
	if err != nil {
		t.Fatalf("json.Marshal: %v", err)
	}
	return b
}
