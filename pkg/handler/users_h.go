package handler

import (
	"github.com/gofiber/fiber/v2"
)

func (h *Handler) getUserMusic(ctx *fiber.Ctx) error {
	var tkn struct {
		GuestID string `json:"gi"`
	}

	if err := ctx.BodyParser(&tkn); err != nil {
		return err
	}

	users, err := h.services.GetUserMusic(tkn.GuestID)

	if err != nil {
		return ctx.JSON(fiber.Map{"error": err.Error()})
	}

	return ctx.JSON(users)
}
