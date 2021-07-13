package handler

import (
	"github.com/gofiber/fiber/v2"
)

func (h *Handler) allUsers(ctx *fiber.Ctx) error {
	users, err := h.services.GetAllUsersDB()

	if err != nil {
		return ctx.JSON(fiber.Map{"error": err.Error()})
	}

	return ctx.JSON(users)
}
