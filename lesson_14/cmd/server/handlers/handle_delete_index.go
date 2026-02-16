package handlers

import (
	"context"
	"time"

	"github.com/gofiber/fiber/v2"
)

type DeleteIndexRequest struct {
	CollectionName string `json:"collection_name"`
	Name           string `json:"name"`
}

func (s *Server) HandleDeleteIndex(c *fiber.Ctx) error {
	var req DeleteIndexRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(OkResponse{Ok: false, Error: "invalid json"})
	}
	if req.CollectionName == "" || req.Name == "" {
		return c.Status(fiber.StatusBadRequest).JSON(OkResponse{Ok: false, Error: "collection_name and name are required"})
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := s.repo.DeleteIndex(ctx, req.CollectionName, req.Name)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(OkResponse{Ok: false, Error: err.Error()})
	}
	return c.JSON(OkResponse{Ok: true})
}
