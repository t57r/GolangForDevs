package main

import (
	"fmt"
	"log"

	"lesson9/internal/documentstore"
)

func strPtr(s string) *string { return &s }

func main() {
	coll := &documentstore.Collection{
		Config: documentstore.CollectionConfig{
			PrimaryKey: "ID",
		},
		Items: make(map[string]*documentstore.Document),
	}

	put(coll, "1", "Apple")
	put(coll, "2", "Bread")
	put(coll, "3", "Chicken")
	put(coll, "4", "Dough")
	put(coll, "5", "Egg")

	fmt.Println("\n=== Documents inserted ===")
	printAll(coll)

	fmt.Println("\n=== Create index on field 'Name' ===")
	if err := coll.CreateIndex("Name"); err != nil {
		log.Fatal(err)
	}

	fmt.Println("\n=== Query Name ASC ===")
	res, err := coll.Query("Name", documentstore.QueryParams{})
	if err != nil {
		log.Fatal(err)
	}
	printDocs(res)

	fmt.Println("\n=== Query Name DESC ===")
	res, err = coll.Query("Name", documentstore.QueryParams{
		Desc: true,
	})
	if err != nil {
		log.Fatal(err)
	}
	printDocs(res)

	fmt.Println("\n=== Query Name BETWEEN 'Bread' AND 'Dough' ===")
	res, err = coll.Query("Name", documentstore.QueryParams{
		MinValue: strPtr("Bread"),
		MaxValue: strPtr("Dough"),
	})
	if err != nil {
		log.Fatal(err)
	}
	printDocs(res)

	fmt.Println("\n=== Update ID=3 Chicken -> Beef ===")
	put(coll, "3", "Beef")

	fmt.Println("\n=== Query Name ASC (after update) ===")
	res, _ = coll.Query("Name", documentstore.QueryParams{})
	printDocs(res)

	fmt.Println("\n=== Delete ID=2 (Bread) ===")
	coll.Delete("2")

	fmt.Println("\n=== Query Name ASC (after delete) ===")
	res, _ = coll.Query("Name", documentstore.QueryParams{})
	printDocs(res)

	fmt.Println("\n=== Delete index on Name ===")
	if err := coll.DeleteIndex("Name"); err != nil {
		log.Fatal(err)
	}

	fmt.Println("\n=== Query after index deletion (should fail) ===")
	_, err = coll.Query("Name", documentstore.QueryParams{})
	fmt.Println("Expected error:", err)
}

func put(c *documentstore.Collection, id, name string) {
	doc := documentstore.Document{
		Fields: map[string]documentstore.DocumentField{
			"ID": {
				Type:  documentstore.DocumentFieldTypeString,
				Value: id,
			},
			"Name": {
				Type:  documentstore.DocumentFieldTypeString,
				Value: name,
			},
		},
	}
	if err := c.Put(doc); err != nil {
		log.Fatal(err)
	}
}

func printAll(c *documentstore.Collection) {
	docs := c.List()
	printDocs(docs)
}

func printDocs(docs []documentstore.Document) {
	for i, d := range docs {
		id := d.Fields["ID"].Value.(string)
		name := d.Fields["Name"].Value.(string)
		fmt.Printf("[%d] ID=%s Name=%s\n", i, id, name)
	}
}
