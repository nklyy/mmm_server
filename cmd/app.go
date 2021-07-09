package cmd

import (
	"fmt"
	"mmm_server/controllers"
	"mmm_server/databases"
	"mmm_server/middleware"
	"mmm_server/repositories"
	"mmm_server/routes"

	"github.com/gofiber/fiber/v2"
)

func Execute() {
	app := fiber.New()

	db, err := databases.MongoDbConnection()
	if err != nil {
		return
	}

	middleware.FiberMiddleware(app)

	repos := repositories.NewRepository(db)
	controll := controllers.NewController(repos)
	route := routes.NewRoute(controll)

	route.Initialroute(app)

	err = app.Listen(":4000")
	if err != nil {
		fmt.Printf("ERROR: %s \n", err)
		return
	}
}

// DNWCIEGKv32vUryK
