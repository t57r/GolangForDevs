package handlers

import (
	"context"
	"time"

	"github.com/gofiber/fiber/v2"
)

type CreateCollectionRequest struct {
	CollectionName string `json:"collection_name"`
}

func (s *Server) HandleCreateCollection(c *fiber.Ctx) error {
	var req CreateCollectionRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(OkResponse{Ok: false, Error: "invalid json"})
	}
	if req.CollectionName == "" {
		return c.Status(fiber.StatusBadRequest).JSON(OkResponse{Ok: false, Error: "collection_name is required"})
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := s.repo.CreateCollection(ctx, req.CollectionName); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(OkResponse{Ok: false, Error: err.Error()})
	}
	return c.JSON(OkResponse{Ok: true})
}
