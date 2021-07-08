package cmd

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"mmm_server/pkg/routes"
)

func Execute() {
	app := fiber.New()

	app.Use(cors.New())

	routes.DeezerRoutes(app)
	routes.NotFoundRoute(app)

	err := app.Listen(":4000")
	if err != nil {
		fmt.Printf("ERROR: %s \n", err)
		return
	}
}
