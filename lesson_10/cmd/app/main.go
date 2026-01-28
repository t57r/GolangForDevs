package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"

	"lesson10/internal/documentstore"
)

func strPtr(s string) *string { return &s }

const (
	collectionName = "concurrent"
	goroutines     = 1000
	keySpace       = 200 // number of distinct IDs being touched
)

func main() {

	store := documentstore.NewStore()

	collAny, err := store.CreateCollection(collectionName, documentstore.CollectionConfig{PrimaryKey: "ID"})
	if err != nil {
		fmt.Printf("CreateCollection error: %v\n", err)
		return
	}

	coll := collAny.(*documentstore.Collection)

	err = coll.CreateIndex("Name")
	if err != nil {
		fmt.Printf("Create Index error: %v\n", err)
	}

	start := time.Now()
	var wg sync.WaitGroup
	wg.Add(goroutines)

	for i := range goroutines {
		go func(worker int) {
			defer wg.Done()

			// Each worker will do a small batch of operations
			ops := 200
			r := rand.New(rand.NewSource(time.Now().UnixNano() + int64(worker)))

			for range ops {
				id := fmt.Sprintf("%d", r.Intn(keySpace))

				switch r.Intn(3) {
				case 0: // Put
					doc := documentstore.Document{
						Fields: map[string]documentstore.DocumentField{
							"ID": {Type: documentstore.DocumentFieldTypeString, Value: id},
							"Name": {
								Type:  documentstore.DocumentFieldTypeString,
								Value: fmt.Sprintf("user-%s-%d", id, worker),
							},
						},
					}
					_ = coll.Put(doc) // ignore errors for stress test

				case 1: // Get
					_, _ = coll.Get(id)

				case 2: // Delete
					_ = coll.Delete(id)
				}
			}
		}(i)
	}

	wg.Wait()

	elapsed := time.Since(start)
	fmt.Printf("Finished %d goroutines in %s\n", goroutines, elapsed)

	docs := coll.List()
	fmt.Printf("Final documents count: %d\n", len(docs))

	res, err := coll.Query("Name", documentstore.QueryParams{MinValue: strPtr("user-")})
	if err == nil {
		fmt.Printf("Query returned %d docs\n", len(res))
	} else {
		fmt.Printf("Query error (ok if index not available): %v\n", err)
	}
}
