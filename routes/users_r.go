package routes

import (
	"mmm_server/controllers"

	"github.com/gofiber/fiber/v2"
)

func UsersRoute(app *fiber.App) {
	api := app.Group("/v1")

	api.Get("/users", controllers.GetAllUsers)
}
