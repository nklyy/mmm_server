package handler

import (
	"github.com/gofiber/fiber/v2"
)

func (h *Handler) spotifyAuthRedirect(ctx *fiber.Ctx) error {
	return ctx.Redirect("https://accounts.spotify.com/authorize?client_id=a45422e6fcc04cc6932840b3372581f5&response_type=code&redirect_uri=http://localhost:4000/v1/spotify/callback&scope=user-read-private%20user-read-email%20user-read-playback-state%20user-modify-playback-state%20user-library-modify%20user-library-read&show_dialog=true")
}

func (h *Handler) spotifyCallback(ctx *fiber.Ctx) error {
	code := ctx.Query("code")

	token := h.services.GetSpotifyAccessToken(code)

	return ctx.Redirect("http://localhost:3000/cf?type=s&code=" + token)
}

func (h *Handler) checkSpotifyAccessToken(ctx *fiber.Ctx) error {
	var cd struct {
		Code string `json:"code"`
	}

	if err := ctx.BodyParser(&cd); err != nil {
		return err
	}

	ok := h.services.CheckSpotifyAccessToken(cd.Code)
	if !ok {
		//return fiber.NewError(fiber.StatusBadRequest, "Invalid token!")
		return ctx.JSON(fiber.Map{"error": "Invalid token!"})
	}

	return ctx.JSON(fiber.Map{"message": "success"})
}
