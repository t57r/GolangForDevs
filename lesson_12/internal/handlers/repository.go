package handlers

import (
	"lesson12/internal/documentstore"
)

type Repository interface {
	GetDocumentStore() *documentstore.Store
	GetCollectionAsCollection(name string) (*documentstore.Collection, error)
}
