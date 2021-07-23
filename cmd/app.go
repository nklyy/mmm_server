package cmd

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"mmm_server/pkg/handler"
	"mmm_server/pkg/repository"
	"mmm_server/pkg/service"
)

func Execute() {
	app := fiber.New()

	// Connection to DB
	db, err := repository.MongoDbConnection()
	if err != nil {
		return
	}

	// Init App Middleware
	app.Use(
		// Add CORS to each route.
		cors.New(),
	)

	// Optional middleware
	app.Use("/v1/ws/deezer/move", func(c *fiber.Ctx) error {
		if c.Get("host") == "localhost:4000" {
			c.Locals("Host", "Localhost:4000")
			return c.Next()
		}
		return c.Status(403).SendString("Request origin not allowed")
	})

	// Init repository, service and handlers
	newRepository := repository.NewRepository(db)
	newService := service.NewService(newRepository)
	newHandler := handler.NewHandler(newService)

	newHandler.InitialRoute(app)

	// NotFound Urls
	app.Use(
		// Anonymous function.
		func(c *fiber.Ctx) error {
			// Return HTTP 404 status and JSON response.
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": true,
				"msg":   "Sorry, endpoint " + "'" + c.OriginalURL() + "'" + " is not found",
			})
		},
	)

	// Starting App
	err = app.Listen(":4000")
	if err != nil {
		fmt.Printf("ERROR: %s \n", err)
		return
	}
}
