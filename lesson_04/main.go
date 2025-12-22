package main

import (
	"fmt"
	"log"

	"lesson4/documentstore"
)

func main() {
	store := documentstore.NewStore()
	cfg := &documentstore.CollectionConfig{PrimaryKey: "id"}

	// Create users
	usersCreated, usersCollection := store.CreateCollection("users", cfg)
	if !usersCreated {
		log.Fatalf("Couldn't create 'users' collection")
	}
	user1 := documentstore.Document{
		Fields: map[string]documentstore.DocumentField{
			"id": {
				Type:  documentstore.DocumentFieldTypeString,
				Value: "user1",
			},
			"name": {
				Type:  documentstore.DocumentFieldTypeString,
				Value: "Taras",
			},
		},
	}
	user2 := documentstore.Document{
		Fields: map[string]documentstore.DocumentField{
			"id": {
				Type:  documentstore.DocumentFieldTypeString,
				Value: "user2",
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
	usersCollection.Put(user1)
	usersCollection.Put(user2)

	// Create products
	productsCreated, productsCollection := store.CreateCollection("products", cfg)
	if !productsCreated {
		log.Fatalf("Couldn't create 'products' collection")
	}
	product1 := documentstore.Document{
		Fields: map[string]documentstore.DocumentField{
			"id": {
				Type:  documentstore.DocumentFieldTypeString,
				Value: "product1",
			},
			"title": {
				Type:  documentstore.DocumentFieldTypeString,
				Value: "Milk",
			},
		},
	}
	product2 := documentstore.Document{
		Fields: map[string]documentstore.DocumentField{
			"id": {
				Type:  documentstore.DocumentFieldTypeString,
				Value: "product2",
			},
			"title": {
				Type:  documentstore.DocumentFieldTypeString,
				Value: "Bread",
			},
		},
	}
	productsCollection.Put(product1)
	productsCollection.Put(product2)

	// Get and display collections
	getAndPrintCollection(store, "users")
	getAndPrintCollection(store, "products")

	deleteByKey(usersCollection, "user2")
	deleteByKey(productsCollection, "product1")

	fmt.Println("\nAfter 'user2' and 'product1' deletion:")
	getAndPrintCollection(store, "users")
	getAndPrintCollection(store, "products")
}

func getAndPrintCollection(store *documentstore.Store, collectionName string) {
	fmt.Printf("\n=== %s\n", collectionName)
	fetchedCollection, exist := store.GetCollection(collectionName)
	if !exist {
		log.Fatalf("Couldn't fetch '%s' from store", collectionName)
	}
	printCollection(fetchedCollection)
}

func printCollection(coll *documentstore.Collection) {
	docs := coll.List()
	fmt.Printf("Total documents: %d\n", len(docs))
	for i, d := range docs {
		fmt.Printf("Document #%d:\n", i+1)
		printDocument(&d)
	}
}

func printDocument(doc *documentstore.Document) {
	for k, f := range doc.Fields {
		fmt.Printf("  %s (%s): %v\n", k, f.Type, f.Value)
	}
}

type Deletable interface {
	Delete(key string) bool
}

func deleteByKey(deletable Deletable, key string) {
	deletable.Delete(key)
}
