package model

import (
	"encoding/json"

	"lesson12/internal/documentstore"
)

type RequestOperation string

const (
	RequestOperationCreateCollection RequestOperation = "create_collection"
	RequestOperationGetCollection    RequestOperation = "get_collection"
	RequestOperationDeleteCollection RequestOperation = "delete_collection"

	RequestOperationPutDocument    RequestOperation = "put_doc"
	RequestOperationGetDocument    RequestOperation = "get_doc"
	RequestOperationDeleteDocument RequestOperation = "delete_doc"
	RequestOperationListDocuments  RequestOperation = "list_docs"

	RequestOperationCreateIndex RequestOperation = "create_index"
	RequestOperationDeleteIndex RequestOperation = "delete_index"
	RequestOperationQuery       RequestOperation = "query"
)

// --- Request/Response ---

type Request struct {
	ID         string           `json:"id"`
	Op         RequestOperation `json:"op"`
	Collection string           `json:"collection,omitempty"`
	Payload    json.RawMessage  `json:"payload,omitempty"`
}

type Response struct {
	ID    string          `json:"id"`
	Ok    bool            `json:"ok"`
	Error string          `json:"error,omitempty"`
	Data  json.RawMessage `json:"data,omitempty"`
}

// --- Payload types ---

type CreateCollectionPayload struct {
	PrimaryKey string `json:"primaryKey"`
}

type DeleteCollectionPayload struct {
	// empty
}

type PutDocPayload struct {
	Doc documentstore.Document `json:"doc"`
}

type GetDocPayload struct {
	Key string `json:"key"`
}

type DeleteDocPayload struct {
	Key string `json:"key"`
}

type CreateIndexPayload struct {
	FieldName string `json:"fieldName"`
}

type DeleteIndexPayload struct {
	FieldName string `json:"fieldName"`
}

type QueryPayload struct {
	FieldName string  `json:"fieldName"`
	Desc      bool    `json:"desc"`
	MinValue  *string `json:"minValue,omitempty"`
	MaxValue  *string `json:"maxValue,omitempty"`
}
