package handlers

import (
	"encoding/json"

	"lesson12/internal/documentstore"
	"lesson12/internal/model"
)

func HandleCreateCollection(repository Repository, req model.Request) model.Response {
	var p model.CreateCollectionPayload
	if err := json.Unmarshal(req.Payload, &p); err != nil {
		return model.BadResponse(req, err)
	}
	coll, err := repository.GetDocumentStore().CreateCollection(req.Collection, documentstore.CollectionConfig{PrimaryKey: p.PrimaryKey})
	if err != nil {
		return model.BadResponse(req, err)
	}
	_ = coll // we return minimal info
	return model.OkResponse(req, map[string]any{"collection": req.Collection})
}

func HandleGetCollection(repository Repository, req model.Request) model.Response {
	_, err := repository.GetDocumentStore().GetCollection(req.Collection)
	if err != nil {
		return model.BadResponse(req, err)
	}
	return model.OkResponse(req, map[string]any{"collection": req.Collection})
}

func HandleDeleteCollection(repository Repository, req model.Request) model.Response {
	if err := repository.GetDocumentStore().DeleteCollection(req.Collection); err != nil {
		return model.BadResponse(req, err)
	}
	return model.OkResponse(req, map[string]any{"deleted": true})
}
