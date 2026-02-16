package handlers

import (
	"context"
	"time"

	"github.com/gofiber/fiber/v2"
)

func (s *Server) HandleListCollections(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	names, err := s.repo.ListCollections(ctx)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(OkResponse{Ok: false, Error: err.Error()})
	}
	return c.JSON(fiber.Map{"ok": true, "collections": names})
}
