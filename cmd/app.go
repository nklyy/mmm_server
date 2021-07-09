package cmd

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"mmm_server/pkg/handler"
	middleware2 "mmm_server/pkg/middleware"
	"mmm_server/pkg/repository"
	"mmm_server/pkg/service"
)

func Execute() {
	app := fiber.New()

	db, err := repository.MongoDbConnection()
	if err != nil {
		return
	}

	middleware2.FiberMiddleware(app)

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

	err = app.Listen(":4000")
	if err != nil {
		fmt.Printf("ERROR: %s \n", err)
		return
	}
}

// DNWCIEGKv32vUryK
