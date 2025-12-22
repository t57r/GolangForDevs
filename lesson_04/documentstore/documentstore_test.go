package documentstore

import "testing"

func TestNewStoreInitializesCollections(t *testing.T) {
	s := NewStore()
	if s == nil {
		t.Fatalf("NewStore() returned nil")
	}
	if s.Collections == nil {
		t.Fatalf("NewStore() must initialize Collections map")
	}
}

func TestCreateCollectionAndGetCollection(t *testing.T) {
	s := NewStore()

	cfg := &CollectionConfig{PrimaryKey: "id"}

	created, collection := s.CreateCollection("users", cfg)
	if !created {
		t.Fatalf("expected first CreateCollection to return true, got false")
	}
	if collection == nil {
		t.Fatalf("expected non-nil Collection")
	}
	if collection.Config == nil || collection.Config.PrimaryKey != "id" {
		t.Fatalf("expected collection config with PrimaryKey 'id', got %#v", collection.Config)
	}
	if collection.Items == nil {
		t.Fatalf("collection.Items must be initialized (non-nil)")
	}

	// Creating the same collection again should fail
	createdAgain, collectionAgain := s.CreateCollection("users", cfg)
	if createdAgain {
		t.Fatalf("expcted second CreateCollection to return false for existing name")
	}
	if collectionAgain != nil {
		t.Fatalf("expected nil collection on duplicate create, got %#v", collectionAgain)
	}

	// GetCollection should find it
	got, ok := s.GetCollection("users")
	if !ok || got == nil {
		t.Fatalf("GetCollection should return existing collection")
	}
	if got != collection {
		t.Fatalf("GetCollection returned different pointer than CreateCollection")
	}

	// GetCollection for non-existing
	if got, ok := s.GetCollection("products"); ok || got != nil {
		t.Fatalf("expected GetCollection for non-existing to return (nil,false), got(%#v,%v)", got, ok)
	}
}

func TestDeleteCollection(t *testing.T) {
	s := NewStore()
	cfg := &CollectionConfig{PrimaryKey: "id"}

	ok, _ := s.CreateCollection("users", cfg)
	if !ok {
		t.Fatalf("failed to create collection in setup")
	}

	// Delete existing
	deleted := s.DeleteCollection("users")
	if !deleted {
		t.Fatalf("expected DeleteCollection to return true for existing collection")
	}

	// Ensure it's gone
	if coll, ok := s.GetCollection("users"); ok || coll != nil {
		t.Fatalf("expected collection 'users' to be removed")
	}

	// Delete non-existing
	deleted = s.DeleteCollection("users")
	if deleted {
		t.Fatalf("expected DeleteCollection to return false for non-existing collection")
	}
}

func newTestCollection(primaryKey string) *Collection {
	return &Collection{
		Config: &CollectionConfig{PrimaryKey: primaryKey},
		Items:  make(map[string]*Document),
	}
}

func TestCollectionPutGetListDelete(t *testing.T) {
	coll := newTestCollection("id")

	doc1 := Document{
		Fields: map[string]DocumentField{
			"id": {
				Type:  DocumentFieldTypeString,
				Value: "user1",
			},
			"name": {
				Type:  DocumentFieldTypeString,
				Value: "Taras",
			},
		},
	}

	doc2 := Document{
		Fields: map[string]DocumentField{
			"id": {
				Type:  DocumentFieldTypeString,
				Value: "user2",
			},
			"age": {
				Type:  DocumentFieldTypeNumber,
				Value: 30,
			},
		},
	}

	coll.Put(doc1)
	coll.Put(doc2)

	if len(coll.Items) != 2 {
		t.Fatalf("expected 2 items after Put, got %d", len(coll.Items))
	}

	// Get existing
	got, ok := coll.Get("user1")
	if !ok || got == nil {
		t.Fatalf("expected Get to find 'user1'")
	}
	if got.Fields["name"].Value != "Taras" {
		t.Fatalf("expected name 'Taras', got %#v", got.Fields["name"].Value)
	}

	// Get non-existing
	got, ok = coll.Get("user404")
	if ok || got != nil {
		t.Fatalf("expected Get non-existing to return (nil,false), got (%#v,%v)", got, ok)
	}

	// List
	list := coll.List()
	if len(list) != 2 {
		t.Fatalf("List: expected 2 docs, got %d", len(list))
	}

	// Delete existing
	deleted := coll.Delete("user1")
	if !deleted {
		t.Fatalf("expected Delete to return true for existing key")
	}
	if _, ok := coll.Get("user1"); ok {
		t.Fatalf("expected 'user1' to be deleted")
	}

	// Delete non-existing
	deleted = coll.Delete("user404")
	if deleted {
		t.Fatalf("expected Delete to return false for non-existing key")
	}
}

func TestCollectionPutMissingPrimaryKey(t *testing.T) {
	coll := newTestCollection("id")

	// No "id" field
	doc := Document{
		Fields: map[string]DocumentField{
			"name": {
				Type:  DocumentFieldTypeString,
				Value: "NoID",
			},
		},
	}

	coll.Put(doc)

	if len(coll.Items) != 0 {
		t.Fatalf("expected no items when primary key field is missing, got %d", len(coll.Items))
	}
}

func TestCollectionPutWrongPrimaryKeyType(t *testing.T) {
	coll := newTestCollection("id")

	// "id" field exists but not string type
	doc := Document{
		Fields: map[string]DocumentField{
			"id": {
				Type:  DocumentFieldTypeNumber, // wrong type
				Value: 123,
			},
		},
	}

	coll.Put(doc)

	if len(coll.Items) != 0 {
		t.Fatalf("expected no items when primary key field has wrong type, got %d", len(coll.Items))
	}
}
