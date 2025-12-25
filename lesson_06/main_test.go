package main

import (
	"os"
	"path/filepath"
	"sort"
	"testing"

	"lesson6/documentstore"
)

func TestStoreDumpRoundTripBytes(t *testing.T) {
	store, err := createProductsStoreExample()
	if err != nil {
		t.Fatalf("createProductsStoreExample() error: %v", err)
	}

	origProducts, err := readProductsFromStore(store)
	if err != nil {
		t.Fatalf("readProductsFromStore(orig) error: %v", err)
	}

	dump, err := store.Dump()
	if err != nil {
		t.Fatalf("store.Dump() error: %v", err)
	}
	if len(dump) == 0 {
		t.Fatalf("store.Dump() returned empty dump")
	}

	restored, err := documentstore.NewStoreFromDump(dump)
	if err != nil {
		t.Fatalf("NewStoreFromDump() error: %v", err)
	}

	restoredProducts, err := readProductsFromStore(restored)
	if err != nil {
		t.Fatalf("readProductsFromStore(restored) error: %v", err)
	}

	assertSameProducts(t, origProducts, restoredProducts)
}

func TestStoreDumpRoundTripFile(t *testing.T) {
	store, err := createProductsStoreExample()
	if err != nil {
		t.Fatalf("createProductsStoreExample() error: %v", err)
	}

	origProducts, err := readProductsFromStore(store)
	if err != nil {
		t.Fatalf("readProductsFromStore(orig) error: %v", err)
	}

	dir := t.TempDir()
	filename := filepath.Join(dir, "my_products.dump")

	if err := store.DumpToFile(filename); err != nil {
		t.Fatalf("DumpToFile(%q) error: %v", filename, err)
	}

	// sanity: file exists + non-empty
	info, err := os.Stat(filename)
	if err != nil {
		t.Fatalf("expected dump file to exist, stat error: %v", err)
	}
	if info.Size() == 0 {
		t.Fatalf("expected dump file to be non-empty")
	}

	restored, err := documentstore.NewStoreFromFile(filename)
	if err != nil {
		t.Fatalf("NewStoreFromFile(%q) error: %v", filename, err)
	}

	restoredProducts, err := readProductsFromStore(restored)
	if err != nil {
		t.Fatalf("readProductsFromStore(restored) error: %v", err)
	}

	assertSameProducts(t, origProducts, restoredProducts)
}

// helper (test-only)
func assertSameProducts(t *testing.T, a, b []Product) {
	t.Helper()

	normalize := func(xs []Product) []Product {
		cp := make([]Product, len(xs))
		copy(cp, xs)
		sort.Slice(cp, func(i, j int) bool {
			// stable ordering by ID then Name
			if cp[i].ID == cp[j].ID {
				return cp[i].Name < cp[j].Name
			}
			return cp[i].ID < cp[j].ID
		})
		return cp
	}

	a2 := normalize(a)
	b2 := normalize(b)

	if len(a2) != len(b2) {
		t.Fatalf("product count mismatch: orig=%d restored=%d\norig=%#v\nrestored=%#v", len(a2), len(b2), a2, b2)
	}

	for i := range a2 {
		if a2[i] != b2[i] {
			t.Fatalf("product mismatch at index %d:\norig=%#v\nrestored=%#v\nall orig=%#v\nall restored=%#v",
				i, a2[i], b2[i], a2, b2)
		}
	}
}
