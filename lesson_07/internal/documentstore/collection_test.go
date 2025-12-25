package documentstore

import (
	"errors"
	"testing"
)

func newTestCollection(pk string) *Collection {
	return &Collection{
		Config: &CollectionConfig{PrimaryKey: pk},
		Items:  make(map[string]*Document),
	}
}

func TestCollectionPut_ConfigNil(t *testing.T) {
	c := &Collection{
		Config: nil,
		Items:  make(map[string]*Document),
	}

	err := c.Put(Document{Fields: map[string]DocumentField{}})
	if !errors.Is(err, ErrConfigNotFound) {
		t.Fatalf("expected ErrConfigNotFound, got %v", err)
	}
}

func TestCollectionPut_InvalidPrimaryKeyField(t *testing.T) {
	c := newTestCollection("ID")

	// Missing primary key
	err := c.Put(Document{
		Fields: map[string]DocumentField{
			"Name": {Type: DocumentFieldTypeString, Value: "Alice"},
		},
	})
	if !errors.Is(err, ErrUnsupportedDocumentField) {
		t.Fatalf("missing pk: expected ErrUnsupportedDocumentField, got %v", err)
	}

	// Wrong type primary key
	err = c.Put(Document{
		Fields: map[string]DocumentField{
			"ID": {Type: DocumentFieldTypeNumber, Value: 123},
		},
	})
	if !errors.Is(err, ErrUnsupportedDocumentField) {
		t.Fatalf("wrong pk type: expected ErrUnsupportedDocumentField, got %v", err)
	}
}

func TestCollectionPutCreateThenUpdate(t *testing.T) {
	c := newTestCollection("ID")

	doc1 := Document{
		Fields: map[string]DocumentField{
			"ID":   {Type: DocumentFieldTypeString, Value: "1"},
			"Name": {Type: DocumentFieldTypeString, Value: "Alice"},
		},
	}
	if err := c.Put(doc1); err != nil {
		t.Fatalf("Put(create) error: %v", err)
	}
	if len(c.Items) != 1 {
		t.Fatalf("expected 1 item, got %d", len(c.Items))
	}

	// Update same ID
	doc2 := Document{
		Fields: map[string]DocumentField{
			"ID":   {Type: DocumentFieldTypeString, Value: "1"},
			"Name": {Type: DocumentFieldTypeString, Value: "AliceUpdated"},
		},
	}
	if err := c.Put(doc2); err != nil {
		t.Fatalf("Put(update) error: %v", err)
	}
	if len(c.Items) != 1 {
		t.Fatalf("expected still 1 item after update, got %d", len(c.Items))
	}

	got, ok := c.Get("1")
	if !ok || got == nil {
		t.Fatalf("expected Get to find key=1")
	}
	name := got.Fields["Name"].Value.(string)
	if name != "AliceUpdated" {
		t.Fatalf("expected updated name, got %q", name)
	}
}

func TestCollectionGet_Delete_List(t *testing.T) {
	c := newTestCollection("ID")

	// Get missing
	if doc, ok := c.Get("missing"); ok || doc != nil {
		t.Fatalf("expected Get(missing) -> (nil,false), got (%v,%v)", doc, ok)
	}

	// Put 2 docs
	_ = c.Put(Document{Fields: map[string]DocumentField{
		"ID": {Type: DocumentFieldTypeString, Value: "1"},
		"X":  {Type: DocumentFieldTypeNumber, Value: int64(10)},
	}})
	_ = c.Put(Document{Fields: map[string]DocumentField{
		"ID": {Type: DocumentFieldTypeString, Value: "2"},
		"X":  {Type: DocumentFieldTypeNumber, Value: int64(20)},
	}})

	// List size
	list := c.List()
	if len(list) != 2 {
		t.Fatalf("expected List() len=2, got %d", len(list))
	}

	// Delete missing
	if ok := c.Delete("missing"); ok {
		t.Fatalf("expected Delete(missing)=false")
	}

	// Delete existing
	if ok := c.Delete("1"); !ok {
		t.Fatalf("expected Delete(1)=true")
	}
	if _, ok := c.Get("1"); ok {
		t.Fatalf("expected key=1 removed")
	}
	if len(c.Items) != 1 {
		t.Fatalf("expected 1 remaining item, got %d", len(c.Items))
	}
}
