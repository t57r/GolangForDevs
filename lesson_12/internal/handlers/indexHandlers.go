package handlers

import (
	"encoding/json"

	"lesson12/internal/documentstore"
	"lesson12/internal/model"
)

func HandleCreateIndex(repository Repository, req model.Request) model.Response {
	var p model.CreateIndexPayload
	if err := json.Unmarshal(req.Payload, &p); err != nil {
		return model.BadResponse(req, err)
	}
	coll, err := repository.GetCollectionAsCollection(req.Collection)
	if err != nil {
		return model.BadResponse(req, err)
	}
	if err := coll.CreateIndex(p.FieldName); err != nil {
		return model.BadResponse(req, err)
	}
	return model.OkResponse(req, map[string]any{"indexed": p.FieldName})
}

func HandleDeleteIndex(repository Repository, req model.Request) model.Response {
	var p model.DeleteIndexPayload
	if err := json.Unmarshal(req.Payload, &p); err != nil {
		return model.BadResponse(req, err)
	}
	coll, err := repository.GetCollectionAsCollection(req.Collection)
	if err != nil {
		return model.BadResponse(req, err)
	}
	if err := coll.DeleteIndex(p.FieldName); err != nil {
		return model.BadResponse(req, err)
	}
	return model.OkResponse(req, map[string]any{"deletedIndex": p.FieldName})
}

func HandleQuery(repository Repository, req model.Request) model.Response {
	var p model.QueryPayload
	if err := json.Unmarshal(req.Payload, &p); err != nil {
		return model.BadResponse(req, err)
	}
	coll, err := repository.GetCollectionAsCollection(req.Collection)
	if err != nil {
		return model.BadResponse(req, err)
	}
	params := documentstore.QueryParams{
		Desc:     p.Desc,
		MinValue: p.MinValue,
		MaxValue: p.MaxValue,
	}
	docs, err := coll.Query(p.FieldName, params)
	if err != nil {
		return model.BadResponse(req, err)
	}
	return model.OkResponse(req, docs)
}
