package handler

import (
	"github.com/gofiber/fiber/v2"
	"strings"
)

func (h *Handler) deezerAuthRedirect(ctx *fiber.Ctx) error {
	m := ctx.Query("m")
	questId := ctx.Query("questId")
	return ctx.Redirect(
		"https://connect.deezer.com/oauth/auth.php?app_id=491682&redirect_uri=http://localhost:4000/v1/deezer/callback&perms=basic_access,manage_library&state="+m+","+questId,
		302)
}

func (h *Handler) deezerCallback(ctx *fiber.Ctx) error {
	code := ctx.Query("code")
	state := ctx.Query("state")
	splitState := strings.Split(state, ",")

	token := h.services.GetDeezerAccessToken(code)

	return ctx.Redirect("http://localhost:3000/cf?type=d&code=" + token + "&m=" + splitState[0] + "&qi=" + splitState[1])
}

func (h *Handler) checkDeezerAccessToken(ctx *fiber.Ctx) error {
	var cd struct {
		Code string `json:"code"`
	}

	if err := ctx.BodyParser(&cd); err != nil {
		return err
	}

	ok := h.services.CheckDeezerAccessToken(cd.Code)
	if !ok {
		//return fiber.NewError(fiber.StatusBadRequest, "Invalid token!")
		return ctx.JSON(fiber.Map{"error": "Invalid token!"})
	}

	return ctx.JSON(fiber.Map{"message": "success"})
}

func (h *Handler) deezerUserMusic(ctx *fiber.Ctx) error {
	var tkn struct {
		Token string `json:"token"`
	}

	if err := ctx.BodyParser(&tkn); err != nil {
		return err
	}

	uMusic := h.services.GetDeezerUserMusic(tkn.Token)

	return ctx.JSON(uMusic)
}
