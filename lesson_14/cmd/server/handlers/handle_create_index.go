package handlers

import (
	"context"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type createIndexRequest struct {
	CollectionName string         `json:"collection_name"`
	Keys           map[string]int `json:"keys"` // 1 ascending, -1 descending
	Unique         bool           `json:"unique,omitempty"`
	Name           string         `json:"name,omitempty"`
}

func (s *Server) HandleCreateIndex(c *fiber.Ctx) error {
	var req createIndexRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(OkResponse{Ok: false, Error: "invalid json"})
	}
	if req.CollectionName == "" || len(req.Keys) == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(OkResponse{Ok: false, Error: "collection_name and keys are required"})
	}

	keys := bson.D{}
	for k, v := range req.Keys {
		if v != 1 && v != -1 {
			return c.Status(fiber.StatusBadRequest).JSON(OkResponse{Ok: false, Error: "index direction must be 1 or -1"})
		}
		keys = append(keys, bson.E{Key: k, Value: v})
	}

	model := mongo.IndexModel{
		Keys: keys,
		Options: options.Index().
			SetUnique(req.Unique),
	}
	if req.Name != "" {
		model.Options.SetName(req.Name)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	name, err := s.db.Collection(req.CollectionName).Indexes().CreateOne(ctx, model)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(OkResponse{Ok: false, Error: err.Error()})
	}
	return c.JSON(fiber.Map{"ok": true, "index_name": name})
}
