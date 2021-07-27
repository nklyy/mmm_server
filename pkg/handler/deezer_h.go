package handler

import (
	"encoding/json"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
	"log"
	"strings"
)

func (h *Handler) deezerAuthRedirect(ctx *fiber.Ctx) error {
	m := ctx.Query("m")
	guestID := ctx.Query("gi")
	return ctx.Redirect(
		"https://connect.deezer.com/oauth/auth.php?app_id=491682&redirect_uri="+h.cfg.DeezerRedirectUrl+"&perms="+h.cfg.DeezerScope+"&state="+m+","+guestID,
		302)
}

func (h *Handler) deezerCallback(ctx *fiber.Ctx) error {
	code := ctx.Query("code")
	state := ctx.Query("state")
	splitState := strings.Split(state, ",")

	// Create Guest User
	if splitState[0] == string('f') {
		findAccessToken, err := h.services.GetDeezerAccessToken(code)
		if err != nil {
			errorMessage, _ := json.Marshal(map[string]string{"error": "Wrong code!"})
			return ctx.Status(400).Send(errorMessage)
		}

		err = h.services.CreateGuestUser(splitState[1], findAccessToken)
		if err != nil {
			errorMessage, _ := json.Marshal(map[string]string{"error": "Something wrong!"})
			return ctx.Status(400).Send(errorMessage)
		}
	}

	if splitState[0] == string('t') {
		moveAccessToken, err := h.services.GetDeezerAccessToken(code)
		if err != nil {
			errorMessage, _ := json.Marshal(map[string]string{"error": "Wrong code!"})
			return ctx.Status(400).Send(errorMessage)
		}

		user, _ := h.services.GetUser(splitState[1])
		user.AccessTokenMove = moveAccessToken

		err = h.services.UpdateGuestUser(splitState[1], user)
		if err != nil {
			errorMessage, _ := json.Marshal(map[string]string{"error": "Something wrong!"})
			return ctx.Status(400).Send(errorMessage)
		}
	}

	return ctx.Redirect(h.cfg.FrontEndUrl + "/cf?type=d&&m=" + splitState[0] + "&gi=" + splitState[1])
}

func (h *Handler) checkDeezerAccessToken(ctx *fiber.Ctx) error {
	var cd struct {
		GuestID string `json:"gi"`
	}

	if err := ctx.BodyParser(&cd); err != nil {
		errorMessage, _ := json.Marshal(map[string]string{"error": "Invalid json!"})
		return ctx.Status(400).Send(errorMessage)
	}

	ok := h.services.CheckDeezerAccessToken(cd.GuestID)
	if !ok {
		errorMessage, _ := json.Marshal(map[string]string{"error": "Invalid token!"})
		return ctx.Status(400).Send(errorMessage)
		//return fiber.NewError(fiber.StatusBadRequest, "Invalid token!")
		//return ctx.JSON(fiber.Map{"error": "Invalid token!"})
	}

	successMessage, _ := json.Marshal(map[string]string{"message": "success"})
	return ctx.Status(200).Send(successMessage)
}

func (h *Handler) deezerUserMusic(ctx *fiber.Ctx) error {
	var tkn struct {
		GuestID string `json:"gi"`
	}

	if err := ctx.BodyParser(&tkn); err != nil {
		errorMessage, _ := json.Marshal(map[string]string{"error": "Invalid body!"})
		return ctx.Status(400).Send(errorMessage)
	}

	user, _ := h.services.GetUser(tkn.GuestID)
	uMusic := h.services.GetDeezerUserMusic(tkn.GuestID)

	// Update Guest User Music
	user.Music = uMusic
	err := h.services.UpdateGuestUser(tkn.GuestID, user)
	if err != nil {
		errorMessage, _ := json.Marshal(map[string]string{"error": "Something wrong!"})
		return ctx.Status(400).Send(errorMessage)
	}

	return ctx.JSON(uMusic)
}

func (h *Handler) moveToDeezer(c *websocket.Conn) {
	var message struct {
		GuestID string `json:"gi"`
	}

	fmt.Println("Remote Address Connected", c.RemoteAddr())
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

	h.services.MoveToDeezer(info.AccessTokenMove, info.Music, c, websocket.TextMessage)
}
