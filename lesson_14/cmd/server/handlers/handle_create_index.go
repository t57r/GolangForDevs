package handlers

import (
	"context"
	"time"

	"github.com/gofiber/fiber/v2"
)

type CreateIndexRequest struct {
	CollectionName string         `json:"collection_name"`
	Keys           map[string]int `json:"keys"` // 1 ascending, -1 descending
	Unique         bool           `json:"unique,omitempty"`
	Name           string         `json:"name,omitempty"`
}

func (s *Server) HandleCreateIndex(c *fiber.Ctx) error {
	var req CreateIndexRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(OkResponse{Ok: false, Error: "invalid json"})
	}
	if req.CollectionName == "" || len(req.Keys) == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(OkResponse{Ok: false, Error: "collection_name and keys are required"})
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	name, err := s.repo.CreateIndex(ctx, req.CollectionName, req.Keys, req.Unique, req.Name)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(OkResponse{Ok: false, Error: err.Error()})
	}
	return c.JSON(fiber.Map{"ok": true, "index_name": name})
}
