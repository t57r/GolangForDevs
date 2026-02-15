package handlers

import (
	"context"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type getDocumentRequest struct {
	CollectionName string         `json:"collection_name"`
	Filter         map[string]any `json:"filter"`
}

func (s *Server) HandleGetDocument(c *fiber.Ctx) error {
	var req getDocumentRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(OkResponse{Ok: false, Error: "invalid json"})
	}
	if req.CollectionName == "" || req.Filter == nil {
		return c.Status(fiber.StatusBadRequest).JSON(OkResponse{Ok: false, Error: "collection_name and filter are required"})
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var out bson.M
	err := s.db.Collection(req.CollectionName).FindOne(ctx, bson.M(req.Filter)).Decode(&out)
	if err == mongo.ErrNoDocuments {
		return c.JSON(fiber.Map{"ok": true, "found": false})
	}
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(OkResponse{Ok: false, Error: err.Error()})
	}
	return c.JSON(fiber.Map{"ok": true, "found": true, "document": out})
}
