package documentstore

import (
	"errors"
	"os"
	"path/filepath"
	"testing"
)

func TestNewStore_InitializesMap(t *testing.T) {
	s := NewStore()
	if s == nil {
		t.Fatalf("NewStore returned nil")
	}
	if s.Collections == nil {
		t.Fatalf("NewStore must init Collections map")
	}
	if len(s.Collections) != 0 {
		t.Fatalf("expected empty store, got %d collections", len(s.Collections))
	}
}

func TestStoreCreateGetDeleteCollection(t *testing.T) {
	s := NewStore()

	// Create with nil config
	if _, err := s.CreateCollection("users", nil); !errors.Is(err, ErrConfigNotFound) {
		t.Fatalf("expected ErrConfigNotFound, got %v", err)
	}

	cfg := &CollectionConfig{PrimaryKey: "ID"}

	// Create OK
	coll, err := s.CreateCollection("users", cfg)
	if err != nil {
		t.Fatalf("CreateCollection error: %v", err)
	}
	if coll == nil {
		t.Fatalf("expected non-nil collection")
	}

	// Duplicate create
	if _, err := s.CreateCollection("users", cfg); !errors.Is(err, ErrCollectionAlreadyExist) {
		t.Fatalf("expected ErrCollectionAlreadyExist, got %v", err)
	}

	// Get existing
	got, err := s.GetCollection("users")
	if err != nil {
		t.Fatalf("GetCollection error: %v", err)
	}
	if got == nil {
		t.Fatalf("GetCollection returned nil")
	}

	// Get missing
	if _, err := s.GetCollection("missing"); !errors.Is(err, ErrCollectionNotFound) {
		t.Fatalf("expected ErrCollectionNotFound, got %v", err)
	}

	// Delete missing
	if err := s.DeleteCollection("missing"); !errors.Is(err, ErrCollectionNotFound) {
		t.Fatalf("expected ErrCollectionNotFound, got %v", err)
	}

	// Delete existing
	if err := s.DeleteCollection("users"); err != nil {
		t.Fatalf("DeleteCollection error: %v", err)
	}
	if _, err := s.GetCollection("users"); !errors.Is(err, ErrCollectionNotFound) {
		t.Fatalf("expected not found after delete, got %v", err)
	}
}

func TestStoreDump_RoundTripBytes(t *testing.T) {
	s := NewStore()
	cfg := &CollectionConfig{PrimaryKey: "ID"}
	collI, err := s.CreateCollection("products", cfg)
	if err != nil {
		t.Fatalf("CreateCollection error: %v", err)
	}
	coll := collI.(*Collection)

	// Add one doc
	doc := Document{
		Fields: map[string]DocumentField{
			"ID":   {Type: DocumentFieldTypeString, Value: "1"},
			"Name": {Type: DocumentFieldTypeString, Value: "Apple"},
		},
	}
	if err := coll.Put(doc); err != nil {
		t.Fatalf("Put error: %v", err)
	}

	dump, err := s.Dump()
	if err != nil {
		t.Fatalf("Dump error: %v", err)
	}
	if len(dump) == 0 {
		t.Fatalf("Dump returned empty bytes")
	}

	s2, err := NewStoreFromDump(dump)
	if err != nil {
		t.Fatalf("NewStoreFromDump error: %v", err)
	}
	if s2 == nil {
		t.Fatalf("NewStoreFromDump returned nil store")
	}

	// Validate restored content
	coll2I, err := s2.GetCollection("products")
	if err != nil {
		t.Fatalf("restored GetCollection error: %v", err)
	}
	coll2 := coll2I.(*Collection)

	got, ok := coll2.Get("1")
	if !ok || got == nil {
		t.Fatalf("expected restored doc key=1")
	}
	if got.Fields["Name"].Value.(string) != "Apple" {
		t.Fatalf("expected Name=Apple, got %#v", got.Fields["Name"].Value)
	}
}

func TestStoreDump_Errors(t *testing.T) {
	// nil receiver
	var s *Store
	if _, err := s.Dump(); err == nil {
		t.Fatalf("expected error for nil store Dump")
	}

	// empty dump
	if _, err := NewStoreFromDump(nil); err == nil {
		t.Fatalf("expected error for empty dump")
	}
}

func TestStoreDumpToFile_RoundTrip(t *testing.T) {
	s := NewStore()
	cfg := &CollectionConfig{PrimaryKey: "ID"}
	collI, err := s.CreateCollection("products", cfg)
	if err != nil {
		t.Fatalf("CreateCollection error: %v", err)
	}
	coll := collI.(*Collection)

	_ = coll.Put(Document{Fields: map[string]DocumentField{
		"ID":   {Type: DocumentFieldTypeString, Value: "1"},
		"Name": {Type: DocumentFieldTypeString, Value: "Apple"},
	}})

	dir := t.TempDir()
	filename := filepath.Join(dir, "store.dump")

	if err := s.DumpToFile(filename); err != nil {
		t.Fatalf("DumpToFile error: %v", err)
	}

	info, err := os.Stat(filename)
	if err != nil {
		t.Fatalf("expected dump file exists, stat error: %v", err)
	}
	if info.Size() == 0 {
		t.Fatalf("expected non-empty dump file")
	}

	s2, err := NewStoreFromFile(filename)
	if err != nil {
		t.Fatalf("NewStoreFromFile error: %v", err)
	}

	coll2I, err := s2.GetCollection("products")
	if err != nil {
		t.Fatalf("restored GetCollection error: %v", err)
	}
	coll2 := coll2I.(*Collection)

	got, ok := coll2.Get("1")
	if !ok {
		t.Fatalf("expected restored doc key=1")
	}
	if got.Fields["Name"].Value.(string) != "Apple" {
		t.Fatalf("expected Name=Apple, got %#v", got.Fields["Name"].Value)
	}
}

func TestStoreDumpToFile_Errors(t *testing.T) {
	s := NewStore()

	if err := s.DumpToFile(""); err == nil {
		t.Fatalf("expected error for empty filename")
	}

	if _, err := NewStoreFromFile(""); err == nil {
		t.Fatalf("expected error for empty filename")
	}

	if _, err := NewStoreFromFile("definitely-does-not-exist.dump"); err == nil {
		t.Fatalf("expected error for missing file")
	}
}
