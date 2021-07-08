package cmd

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"mmm_server/pkg/middleware"
	"mmm_server/pkg/routes"
)

func Execute() {
	app := fiber.New()

	middleware.FiberMiddleware(app)

	routes.DeezerRoutes(app)
	routes.NotFoundRoute(app)

	err := app.Listen(":4000")
	if err != nil {
		fmt.Printf("ERROR: %s \n", err)
		return
	}
}
