package handler

import (
	"github.com/gofiber/fiber/v2"
)

func (h *Handler) allUsers(ctx *fiber.Ctx) error {
	users, err := h.services.GetAllUsers()

	if err != nil {
		return err
	}

	return ctx.JSON(users)
}
