package cmd

import (
	"fmt"
	"mmm_server/databases"
	"mmm_server/middleware"
	"mmm_server/routes"

	"github.com/gofiber/fiber/v2"
)

func Execute() {
	app := fiber.New()

	middleware.FiberMiddleware(app)

	routes.UsersRoute(app)
	routes.DeezerRoute(app)
	routes.NotFoundRoute(app)

	_, err := databases.MongoDbConnection()
	if err != nil {
		return
	}

	err = app.Listen(":4000")
	if err != nil {
		fmt.Printf("ERROR: %s \n", err)
		return
	}
}

// DNWCIEGKv32vUryK
