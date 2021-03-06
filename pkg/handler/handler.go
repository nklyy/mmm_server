package handler

import (
	"github.com/gofiber/websocket/v2"
	"mmm_server/config"
	"mmm_server/pkg/service"

	"github.com/gofiber/fiber/v2"
)

type Handler struct {
	services *service.Service
	cfg      *config.Configurations
}

func NewHandler(services *service.Service, cfg *config.Configurations) *Handler {
	return &Handler{
		services: services,
		cfg:      cfg,
	}
}

func (h *Handler) InitialRoute(route fiber.Router) {
	v1 := route.Group("/v1")

	{
		// Deezer
		v1.Get("/deezer", h.deezerAuthRedirect)
		v1.Get("/deezer/callback", h.deezerCallback)
		v1.Post("/deezer/checkT", h.checkDeezerAccessToken)
		v1.Post("/deezer/userMusic", h.deezerUserMusic)
		v1.Get("/ws/deezer/move", websocket.New(h.moveToDeezer))

		// Spotify
		v1.Get("/spotify", h.spotifyAuthRedirect)
		v1.Get("/spotify/callback", h.spotifyCallback)
		v1.Post("/spotify/checkT", h.checkSpotifyAccessToken)
		v1.Post("/spotify/userMusic", h.spotifyUserMusic)
		v1.Get("/ws/spotify/move", websocket.New(h.moveToSpotify))

		// User
		v1.Get("/userMusic", h.getUserMusic)
	}
}
