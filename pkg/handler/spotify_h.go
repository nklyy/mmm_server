package handler

import (
	"github.com/gofiber/fiber/v2"
	"net/url"
	"strings"
)

func (h *Handler) spotifyAuthRedirect(ctx *fiber.Ctx) error {
	m := ctx.Query("m")
	guestID := ctx.Query("guestID")
	scope := url.PathEscape("user-read-private user-read-email user-read-playback-state user-modify-playback-state user-library-modify user-library-read")
	r := url.PathEscape("http://localhost:4000/v1/spotify/callback")

	return ctx.Redirect("https://accounts.spotify.com/authorize?response_type=code&client_id=6b990a58d275455da234d248fda89722&scope=" + scope + "&redirect_uri=" + r + "&show_dialog=true&state=" + m + "," + guestID)
}

func (h *Handler) spotifyCallback(ctx *fiber.Ctx) error {
	code := ctx.Query("code")
	state := ctx.Query("state")
	splitState := strings.Split(state, ",")

	token := h.services.GetSpotifyAccessToken(code)

	/*
	 TODO Save token into database, check all code, make correct handle error.
	*/

	// Create Guest User
	h.services.CreateGuestUser(splitState[1])

	return ctx.Redirect("http://localhost:3000/cf?type=s&code=" + token + "&m=" + splitState[0] + "&gi=" + splitState[1])
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

func (h *Handler) spotifyUserMusic(ctx *fiber.Ctx) error {
	var tkn struct {
		Code    string `json:"code"`
		GuestID string `json:"gi"`
	}

	if err := ctx.BodyParser(&tkn); err != nil {
		return err
	}

	uMusic := h.services.GetSpotifyUserMusic(tkn.Code)

	// Update Guest User Music
	h.services.UpdateGuestUser(tkn.GuestID, uMusic)

	return ctx.JSON(uMusic)
}

func (h *Handler) moveToSpotify(ctx *fiber.Ctx) error {
	var tkn struct {
		Code    string `json:"code"`
		GuestID string `json:"gi"`
	}

	if err := ctx.BodyParser(&tkn); err != nil {
		return err
	}

	// Get Guest User Music
	music, err := h.services.GetUserMusic(tkn.GuestID)
	if err != nil {
		return err
	}

	h.services.MoveToSpotify(music, tkn.Code)

	return ctx.JSON("")
}
