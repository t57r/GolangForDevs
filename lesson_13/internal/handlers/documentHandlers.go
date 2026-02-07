package handlers

import (
	"encoding/json"

	"lesson12/internal/model"
)

func HandlePutDocument(repository Repository, req model.Request) model.Response {
	var p model.PutDocPayload
	if err := json.Unmarshal(req.Payload, &p); err != nil {
		return model.BadResponse(req, err)
	}
	coll, err := repository.GetCollectionAsCollection(req.Collection)
	if err != nil {
		return model.BadResponse(req, err)
	}
	if err := coll.Put(p.Doc); err != nil {
		return model.BadResponse(req, err)
	}
	return model.OkResponse(req, map[string]any{"stored": true})
}

func HandleGetDocument(repository Repository, req model.Request) model.Response {
	var p model.GetDocPayload
	if err := json.Unmarshal(req.Payload, &p); err != nil {
		return model.BadResponse(req, err)
	}
	coll, err := repository.GetCollectionAsCollection(req.Collection)
	if err != nil {
		return model.BadResponse(req, err)
	}
	doc, found := coll.Get(p.Key)
	if !found {
		return model.Response{ID: req.ID, Ok: false, Error: "document not found"}
	}
	return model.OkResponse(req, doc)
}

func HandleDeleteDocument(repository Repository, req model.Request) model.Response {
	var p model.DeleteDocPayload
	if err := json.Unmarshal(req.Payload, &p); err != nil {
		return model.BadResponse(req, err)
	}
	coll, err := repository.GetCollectionAsCollection(req.Collection)
	if err != nil {
		return model.BadResponse(req, err)
	}
	deleted := coll.Delete(p.Key)
	if !deleted {
		return model.Response{ID: req.ID, Ok: false, Error: "document not found"}
	}
	return model.OkResponse(req, map[string]any{"deleted": true})
}

func HandleListDocuments(repository Repository, req model.Request) model.Response {
	coll, err := repository.GetCollectionAsCollection(req.Collection)
	if err != nil {
		return model.BadResponse(req, err)
	}
	docs := coll.List()
	return model.OkResponse(req, docs)
}
