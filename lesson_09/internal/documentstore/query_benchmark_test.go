package documentstore

import (
	"fmt"
	"io"
	"log/slog"
	"testing"

	"github.com/google/btree"
)

// helper: pointer to string
func strPtr(s string) *string { return &s }

// Build a collection with n docs. Each doc has:
// - primary key "ID" (string)
// - indexed field "X" (string)
// X values are lexicographically sortable and unique-ish.
func buildIndexedCollection(n int) (*Collection, []string, error) {
	c := &Collection{
		Config:  CollectionConfig{PrimaryKey: "ID"},
		Items:   make(map[string]*Document, n),
		indexes: make(map[string]*btree.BTree),
	}

	values := make([]string, 0, n)

	for i := range n {
		id := fmt.Sprintf("%d", i)
		x := fmt.Sprintf("%08d", i) // zero-padded for stable lexicographic ordering

		doc := Document{
			Fields: map[string]DocumentField{
				"ID": {Type: DocumentFieldTypeString, Value: id},
				"X":  {Type: DocumentFieldTypeString, Value: x},
			},
		}
		if err := c.Put(doc); err != nil {
			return nil, nil, err
		}
		values = append(values, x)
	}

	if err := c.CreateIndex("X"); err != nil {
		return nil, nil, err
	}

	return c, values, nil
}

func BenchmarkQuery_OlogN_ExactMatch(b *testing.B) {
	logger = slog.New(slog.NewTextHandler(io.Discard, nil))
	defer func() {
		logger = slog.Default()
	}()

	// Sizes that make it easy to see scaling when doubling.
	sizes := []int{1_000, 2_000, 4_000, 8_000, 16_000, 32_000, 64_000, 128_000}

	for _, n := range sizes {
		b.Run(fmt.Sprintf("n=%d", n), func(b *testing.B) {
			c, values, err := buildIndexedCollection(n)
			if err != nil {
				b.Fatalf("setup failed: %v", err)
			}

			v := values[n/2]
			params := QueryParams{
				Desc:     false,
				MinValue: strPtr(v),
				MaxValue: strPtr(v),
			}

			b.ReportAllocs()
			b.ResetTimer()

			for i := 0; i < b.N; i++ {
				res, err := c.Query("X", params)
				if err != nil {
					b.Fatalf("Query failed: %v", err)
				}
				// Expect 1 hit (exact match) if X unique
				if len(res) != 1 {
					b.Fatalf("expected 1 result, got %d", len(res))
				}
			}

			b.StopTimer()
		})
	}
}
