package handler

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
)

func (h *Handler) deezerAuthRedirect(ctx *fiber.Ctx) error {
	return ctx.Redirect("https://connect.deezer.com/oauth/auth.php?app_id=491682&redirect_uri=http://localhost:4000/v1/deezer/callback&perms=basic_access,manage_library", 302)
}

func (h *Handler) deezerCallback(ctx *fiber.Ctx) error {
	code := ctx.Query("code")

	token := h.services.GetDeezerAccessToken(code)

	return ctx.Redirect("http://localhost:3000/cf?code=" + token)
}

func (h *Handler) checkAccessToken(ctx *fiber.Ctx) error {
	var cd struct {
		Code string `json:"code"`
	}

	if err := ctx.BodyParser(&cd); err != nil {
		return err
	}

	ok := h.services.CheckAccessToken(cd.Code)
	if !ok {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid token!")
	}

	return ctx.JSON(fiber.Map{"message": "success"})
}

func (h *Handler) deezerUserMusic(ctx *fiber.Ctx) error {
	token := ctx.Body()
	fmt.Println(token)

	//h.services.GetDeezerUserMusic(token)

	return ctx.JSON("")
}
