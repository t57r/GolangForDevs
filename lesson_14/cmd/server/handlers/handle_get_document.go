package handlers

import (
	"context"
	"time"

	"github.com/gofiber/fiber/v2"
)

type GetDocumentRequest struct {
	CollectionName string         `json:"collection_name"`
	Filter         map[string]any `json:"filter"`
}

func (s *Server) HandleGetDocument(c *fiber.Ctx) error {
	var req GetDocumentRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(OkResponse{Ok: false, Error: "invalid json"})
	}
	if req.CollectionName == "" || req.Filter == nil {
		return c.Status(fiber.StatusBadRequest).JSON(OkResponse{Ok: false, Error: "collection_name and filter are required"})
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	res, err := s.repo.GetDocument(ctx, req.CollectionName, req.Filter)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(OkResponse{Ok: false, Error: err.Error()})
	}
	return c.JSON(fiber.Map(res))
}
