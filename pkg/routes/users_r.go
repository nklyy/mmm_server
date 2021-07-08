package routes

import (
	"github.com/gofiber/fiber/v2"
	"mmm_server/app/controllers"
)

func UsersRoute(app *fiber.App) {
	api := app.Group("/v1")

	api.Get("/users", controllers.GetAllUsers)
}
