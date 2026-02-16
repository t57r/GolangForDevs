package handlers

import (
	"context"
	"time"

	"github.com/gofiber/fiber/v2"
)

type PutDocumentRequest struct {
	CollectionName string         `json:"collection_name"`
	Document       map[string]any `json:"document"`
}

func (s *Server) HandlePutDocument(c *fiber.Ctx) error {
	var req PutDocumentRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(OkResponse{Ok: false, Error: "invalid json"})
	}
	if req.CollectionName == "" || req.Document == nil {
		return c.Status(fiber.StatusBadRequest).JSON(OkResponse{Ok: false, Error: "collection_name and document are required"})
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := s.repo.PutDocument(ctx, req.CollectionName, req.Document)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(OkResponse{Ok: false, Error: err.Error()})
	}
	return c.JSON(OkResponse{Ok: true})
}
