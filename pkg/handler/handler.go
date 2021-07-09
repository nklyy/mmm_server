package handler

import (
	"mmm_server/pkg/service"

	"github.com/gofiber/fiber/v2"
)

type Handler struct {
	services *service.Service
}

func NewHandler(services *service.Service) *Handler {
	return &Handler{services: services}
}

func (h *Handler) InitialRoute(route fiber.Router) {
	user := route.Group("/user")

	{
		user.Get("/", h.allUsers)
	}
}
