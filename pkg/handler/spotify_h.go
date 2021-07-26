package handler

import (
	"encoding/json"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
	"log"
	"net/url"
	"strings"
)

func (h *Handler) spotifyAuthRedirect(ctx *fiber.Ctx) error {
	m := ctx.Query("m")
	guestID := ctx.Query("gi")
	scope := url.PathEscape(h.cfg.SpotifyScope)
	r := url.PathEscape(h.cfg.SpotifyRedirectUrl)

	return ctx.Redirect("https://accounts.spotify.com/authorize?response_type=code&client_id=" + h.cfg.SpotifyClientKey + "&scope=" + scope + "&redirect_uri=" + r + "&show_dialog=true&state=" + m + "," + guestID)
}

func (h *Handler) spotifyCallback(ctx *fiber.Ctx) error {
	code := ctx.Query("code")
	state := ctx.Query("state")
	splitState := strings.Split(state, ",")

	// Create Guest User
	if splitState[0] == string('f') {
		findAccessToken := h.services.GetSpotifyAccessToken(code)

		err := h.services.CreateGuestUser(splitState[1], findAccessToken)
		if err != nil {
			errorMessage, _ := json.Marshal(map[string]string{"error": "Something wrong!"})
			return ctx.Status(400).Send(errorMessage)
		}
	}

	if splitState[0] == string('t') {
		accessToken := h.services.GetSpotifyAccessToken(code)

		user, _ := h.services.GetUser(splitState[1])
		user.AccessTokenMove = accessToken

		err := h.services.UpdateGuestUser(splitState[1], user)
		if err != nil {
			errorMessage, _ := json.Marshal(map[string]string{"error": "Something wrong!"})
			return ctx.Status(400).Send(errorMessage)
		}
	}

	return ctx.Redirect(h.cfg.FrontEndUrl + "/cf?type=s&m=" + splitState[0] + "&gi=" + splitState[1])
}

func (h *Handler) checkSpotifyAccessToken(ctx *fiber.Ctx) error {
	var cd struct {
		GuestID string `json:"gi"`
	}

	if err := ctx.BodyParser(&cd); err != nil {
		return err
	}

	ok := h.services.CheckSpotifyAccessToken(cd.GuestID)
	if !ok {
		errorMessage, _ := json.Marshal(map[string]string{"error": "Invalid token!"})
		return ctx.Status(400).Send(errorMessage)
		//return fiber.NewError(fiber.StatusBadRequest, "Invalid token!")
		//return ctx.JSON(fiber.Map{"error": "Invalid token!"})
	}

	successMessage, _ := json.Marshal(map[string]string{"message": "success"})
	return ctx.Status(200).Send(successMessage)
}

func (h *Handler) spotifyUserMusic(ctx *fiber.Ctx) error {
	var tkn struct {
		GuestID string `json:"gi"`
	}

	if err := ctx.BodyParser(&tkn); err != nil {
		errorMessage, _ := json.Marshal(map[string]string{"error": "Invalid token!"})
		return ctx.Status(400).Send(errorMessage)
	}

	user, _ := h.services.GetUser(tkn.GuestID)
	uMusic := h.services.GetSpotifyUserMusic(tkn.GuestID)

	// Update Guest User Music
	user.Music = uMusic
	err := h.services.UpdateGuestUser(tkn.GuestID, user)
	if err != nil {
		errorMessage, _ := json.Marshal(map[string]string{"error": "Something wrong!"})
		return ctx.Status(400).Send(errorMessage)
	}

	return ctx.JSON(uMusic)
}

func (h *Handler) moveToSpotify(c *websocket.Conn) {
	var message struct {
		GuestID string `json:"gi"`
	}

	fmt.Println(c.Locals("Host")) // "Localhost:3000"

	fmt.Println("Remote Address", c.RemoteAddr())
	_, msg, err := c.ReadMessage()
	if err != nil {
		log.Println("read:", err)
		return
	}

	err = json.Unmarshal(msg, &message)
	if err != nil {
		fmt.Println(err)
		return
	}

	//fmt.Println("Message", message.GuestID)
	// Get Guest User Music
	info, err := h.services.GetUser(message.GuestID)
	if err != nil {
		return
	}

	//err = c.WriteMessage(mt, []byte(strconv.Itoa(len(info.Music))))

	h.services.MoveToSpotify(info.AccessTokenMove, info.Music, c, websocket.TextMessage)
}
