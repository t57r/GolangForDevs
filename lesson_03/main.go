package main

import (
	"fmt"

	"lesson3/documentstore"
)

func main() {
	doc1 := &documentstore.Document{
		Fields: map[string]documentstore.DocumentField{
			"key": {
				Type:  documentstore.DocumentFieldTypeString,
				Value: "user1",
			},
			"name": {
				Type:  documentstore.DocumentFieldTypeString,
				Value: "John",
			},
			"age": {
				Type:  documentstore.DocumentFieldTypeNumber,
				Value: 30,
			},
			"tags": {
				Type:  documentstore.DocumentFieldTypeArray,
				Value: [...]string{"tag1", "tag2", "tag3"},
			},
		},
	}

	type Coord struct {
		lat, lng float64
	}
	location := Coord{lat: 43.05, lng: 52.85}

	doc2 := &documentstore.Document{
		Fields: map[string]documentstore.DocumentField{
			"key": {
				Type:  documentstore.DocumentFieldTypeString,
				Value: "user2",
			},
			"name": {
				Type:  documentstore.DocumentFieldTypeString,
				Value: "Tom",
			},
			"active": {
				Type:  documentstore.DocumentFieldTypeBool,
				Value: true,
			},
			"location": {
				Type:  documentstore.DocumentFieldTypeObject,
				Value: location,
			},
		},
	}

	// Put
	documentstore.Put(doc1)
	documentstore.Put(doc2)

	fmt.Println("After Put:")
	printAll()

	// Get
	if got, ok := documentstore.Get("user1"); ok {
		fmt.Println("\nGet(\"user1\") found:")
		printDocument(got)
	} else {
		fmt.Println("\nGet(\"user1\") not found")
	}

	if _, ok := documentstore.Get("user404"); !ok {
		fmt.Println("\nGet(\"user404\") not found (as expected)")
	}

	// Delete
	deleted := documentstore.Delete("user1")
	fmt.Println("\nDelete(\"user1\") ->", deleted)

	deleted = documentstore.Delete("user404")
	fmt.Println("Delete(\"user404\") ->", deleted)

	fmt.Println("\nAfter Delete:")
	printAll()
}

func printAll() {
	docs := documentstore.List()
	fmt.Printf("Total documents: %d\n", len(docs))
	for i, d := range docs {
		fmt.Printf("Document #%d:\n", i+1)
		printDocument(d)
	}
}

func printDocument(doc *documentstore.Document) {
	for k, f := range doc.Fields {
		fmt.Printf("  %s (%s): %v\n", k, f.Type, f.Value)
	}
}
