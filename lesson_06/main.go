package main

import (
	"fmt"
	"strconv"

	"lesson6/documentstore"
)

const PRODUCT_COLLECTION_NAME = "products"

type Product struct {
	ID   string
	Name string
}

func main() {
	store, err := createProductsStoreExample()
	if err != nil {
		fmt.Printf("Error creating products store: %v\n", err)
		return
	}

	fmt.Println("=== Created products:")
	displayStoreProducts(store)

	testByteDumpRoundTrip(store)
	testFileDumpRoundTrip(store)

}

func createProductsStoreExample() (*documentstore.Store, error) {
	store := documentstore.NewStore()
	config := &documentstore.CollectionConfig{PrimaryKey: "ID"}

	productCollection, err := store.CreateCollection(PRODUCT_COLLECTION_NAME, config)
	if err != nil {
		return nil, err
	}

	productNames := []string{"Apple", "Bread", "Chicken", "Dough", "Egg"}
	for index := range 5 {
		id := strconv.Itoa(index)
		newProductDoc, err := createProductDocument(id, productNames[index])
		if err != nil {
			return nil, err
		}
		productCollection.Put(*newProductDoc)
	}

	return store, nil
}

func createProductDocument(id string, name string) (*documentstore.Document, error) {
	newProduct := Product{
		ID:   id,
		Name: name,
	}
	newProductDocument, err := documentstore.MarshalDocument(&newProduct)
	if err != nil {
		return nil, err
	}
	return newProductDocument, nil
}

func readProductsFromStore(store *documentstore.Store) ([]Product, error) {
	coll, err := store.GetCollection(PRODUCT_COLLECTION_NAME)
	if err != nil {
		return nil, err
	}

	docs := coll.List()
	products := make([]Product, 0, len(docs))
	for _, doc := range docs {
		var p Product
		if err := documentstore.UnmarshalDocument(&doc, &p); err != nil {
			return nil, err
		}
		products = append(products, p)
	}
	return products, nil
}

func displayStoreProducts(store *documentstore.Store) {
	products, err := readProductsFromStore(store)
	if err != nil {
		fmt.Printf("Error reading products from store: %v\n", err)
		return
	}

	for index, product := range products {
		fmt.Printf("Product[%d] => %v\n", index, product)
	}
}

func testByteDumpRoundTrip(store *documentstore.Store) {
	fmt.Println("\nTest dumb by bytes")
	storeDump, err := store.Dump()
	if err != nil {
		fmt.Printf("Error dumping the store: %v\n", err)
		return
	}
	fmt.Printf("\nSuccessfully dumped the store into %d bytes\n", len(storeDump))

	restoredStore, err := documentstore.NewStoreFromDump(storeDump)
	if err != nil {
		fmt.Printf("Error restoring from the dump: %v\n", err)
		return
	}

	fmt.Println("\n=== Restored products from the dump:")
	displayStoreProducts(restoredStore)
}

func testFileDumpRoundTrip(store *documentstore.Store) {
	fmt.Println("\nTest dumb by file")

	fileName := "my_products"

	err := store.DumpToFile(fileName)
	if err != nil {
		fmt.Printf("Error dumping to file: %v\n", err)
		return
	}
	fmt.Printf("\nSuccessfully dumped the store into '%s' file\n", fileName)

	restoredStore, err := documentstore.NewStoreFromFile(fileName)
	if err != nil {
		fmt.Printf("Error restoring from the dump: %v\n", err)
		return
	}

	fmt.Println("\n=== Restored products from the file:")
	displayStoreProducts(restoredStore)
}
