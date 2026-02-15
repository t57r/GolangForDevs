package handlers

import (
	"context"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type listDocumentsRequest struct {
	CollectionName string         `json:"collection_name"`
	Filter         map[string]any `json:"filter,omitempty"`
	Limit          int64          `json:"limit,omitempty"`
	Skip           int64          `json:"skip,omitempty"`
}

func (s *Server) HandleListDocuments(c *fiber.Ctx) error {
	var req listDocumentsRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(OkResponse{Ok: false, Error: "invalid json"})
	}
	if req.CollectionName == "" {
		return c.Status(fiber.StatusBadRequest).JSON(OkResponse{Ok: false, Error: "collection_name is required"})
	}
	if req.Limit == 0 {
		req.Limit = 100
	}
	if req.Filter == nil {
		req.Filter = map[string]any{}
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	opts := options.Find().SetLimit(req.Limit).SetSkip(req.Skip)
	cur, err := s.db.Collection(req.CollectionName).Find(ctx, bson.M(req.Filter), opts)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(OkResponse{Ok: false, Error: err.Error()})
	}
	defer cur.Close(ctx)

	var docs []bson.M
	if err := cur.All(ctx, &docs); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(OkResponse{Ok: false, Error: err.Error()})
	}

	return c.JSON(fiber.Map{"ok": true, "documents": docs})
}
