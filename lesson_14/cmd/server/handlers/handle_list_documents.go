package handlers

import (
	"context"
	"time"

	"github.com/gofiber/fiber/v2"
)

type ListDocumentsRequest struct {
	CollectionName string         `json:"collection_name"`
	Filter         map[string]any `json:"filter,omitempty"`
	Limit          int64          `json:"limit,omitempty"`
	Skip           int64          `json:"skip,omitempty"`
}

func (s *Server) HandleListDocuments(c *fiber.Ctx) error {
	var req ListDocumentsRequest
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

	docs, err := s.repo.ListDocuments(ctx, req.CollectionName, req.Filter, req.Limit, req.Skip)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(OkResponse{Ok: false, Error: err.Error()})
	}
	return c.JSON(fiber.Map{"ok": true, "documents": docs})
}
