package main

import (
	"project-mgmt-go/config"
	"project-mgmt-go/db"
	"project-mgmt-go/router"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func main() {
	app := fiber.New()
	app.Use(logger.New())
	app.Use(cors.New())

	db.ConnectDB()

	router.SetUpRoutes(app)

	app.Listen(config.Env("PORT"))
}
