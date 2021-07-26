package cmd

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/spf13/viper"
	"log"
	"mmm_server/config"
	"mmm_server/pkg/handler"
	"mmm_server/pkg/repository"
	"mmm_server/pkg/service"
)

func Execute() {
	cfg, err := initConfig()
	if err != nil {
		log.Fatalf("error initializing config: %s", err.Error())
	}

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

	app.Use("/v1/ws/spotify/move", func(c *fiber.Ctx) error {
		if c.Get("host") == "localhost:4000" {
			c.Locals("Host", "Localhost:4000")
			return c.Next()
		}
		return c.Status(403).SendString("Request origin not allowed")
	})

	// Init repository, service and handlers
	newRepository := repository.NewRepository(db)
	newService := service.NewService(newRepository, cfg)
	newHandler := handler.NewHandler(newService, cfg)

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
	err = app.Listen(cfg.PORT)
	if err != nil {
		fmt.Printf("ERROR: %s \n", err)
		return
	}
}

func initConfig() (*config.Configurations, error) {
	viper.AddConfigPath("config")

	viper.SetConfigName("app")

	viper.SetConfigType("env")

	viper.AutomaticEnv()

	err := viper.ReadInConfig()
	if err != nil {
		return nil, err
	}

	var configuration config.Configurations
	err = viper.Unmarshal(&configuration)
	if err != nil {
		fmt.Printf("Unable to decode into struct, %v", err)
	}

	return &configuration, nil
}
