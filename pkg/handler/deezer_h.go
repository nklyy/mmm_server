package handler

import (
	"github.com/gofiber/fiber/v2"
	"strings"
)

func (h *Handler) deezerAuthRedirect(ctx *fiber.Ctx) error {
	m := ctx.Query("m")
	guestID := ctx.Query("gi")
	return ctx.Redirect(
		"https://connect.deezer.com/oauth/auth.php?app_id=491682&redirect_uri=http://localhost:4000/v1/deezer/callback&perms=basic_access,manage_library&state="+m+","+guestID,
		302)
}

func (h *Handler) deezerCallback(ctx *fiber.Ctx) error {
	code := ctx.Query("code")
	state := ctx.Query("state")
	splitState := strings.Split(state, ",")

	// Create Guest User
	if splitState[0] == string('f') {
		findAccessToken := h.services.GetDeezerAccessToken(code)

		h.services.CreateGuestUser(splitState[1], findAccessToken)
	}

	if splitState[0] == string('t') {
		accessToken := h.services.GetDeezerAccessToken(code)

		user, _ := h.services.GetUser(splitState[1])
		user.AccessTokenMove = accessToken

		h.services.UpdateGuestUser(splitState[1], user)
	}

	return ctx.Redirect("http://localhost:3000/cf?type=d&&m=" + splitState[0] + "&gi=" + splitState[1])
}

func (h *Handler) checkDeezerAccessToken(ctx *fiber.Ctx) error {
	var cd struct {
		GuestID string `json:"gi"`
	}

	if err := ctx.BodyParser(&cd); err != nil {
		return err
	}

	ok := h.services.CheckDeezerAccessToken(cd.GuestID)
	if !ok {
		//return fiber.NewError(fiber.StatusBadRequest, "Invalid token!")
		return ctx.JSON(fiber.Map{"error": "Invalid token!"})
	}

	return ctx.JSON(fiber.Map{"message": "success"})
}

func (h *Handler) deezerUserMusic(ctx *fiber.Ctx) error {
	var tkn struct {
		GuestID string `json:"gi"`
	}

	if err := ctx.BodyParser(&tkn); err != nil {
		return err
	}

	user, _ := h.services.GetUser(tkn.GuestID)
	uMusic := h.services.GetDeezerUserMusic(tkn.GuestID)

	// Update Guest User Music
	user.Music = uMusic
	h.services.UpdateGuestUser(tkn.GuestID, user)

	return ctx.JSON(uMusic)
}

func (h *Handler) moveToDeezer(ctx *fiber.Ctx) error {
	var tkn struct {
		GuestID string `json:"gi"`
	}

	if err := ctx.BodyParser(&tkn); err != nil {
		return err
	}

	// Get Guest User Music
	info, err := h.services.GetUser(tkn.GuestID)
	if err != nil {
		return err
	}

	notFoundM := h.services.MoveToDeezer(info.AccessTokenMove, info.Music)

	return ctx.JSON(notFoundM)
}
