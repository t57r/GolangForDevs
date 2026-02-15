package handlers

import (
	"context"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
)

type deleteDocumentRequest struct {
	CollectionName string         `json:"collection_name"`
	Filter         map[string]any `json:"filter"`
}

func (s *Server) HandleDeleteDocument(c *fiber.Ctx) error {
	var req deleteDocumentRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(OkResponse{Ok: false, Error: "invalid json"})
	}
	if req.CollectionName == "" || req.Filter == nil {
		return c.Status(fiber.StatusBadRequest).JSON(OkResponse{Ok: false, Error: "collection_name and filter are required"})
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	res, err := s.db.Collection(req.CollectionName).DeleteOne(ctx, bson.M(req.Filter))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(OkResponse{Ok: false, Error: err.Error()})
	}
	return c.JSON(fiber.Map{"ok": true, "deleted": res.DeletedCount})
}
