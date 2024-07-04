package router

import (
	"net/http"
	"project-mgmt-go/controller"

	"github.com/gofiber/fiber/v2"
)

func SetUpRoutes(app *fiber.App) {
	api := app.Group("/api")

	api.Get("/", func(c *fiber.Ctx) error {
		return c.Status(http.StatusOK).JSON(fiber.Map{"message": "Hello, World!"})
	})

	api.Route("/projects", func(projects fiber.Router) {
		projects.Get("/", controller.GetProjects)
		projects.Get("/:projectId", controller.GetProject)
		projects.Post("/", controller.CreateProject)
		projects.Put("/:projectId", controller.UpdateProject)
		projects.Delete("/:projectId", controller.DeleteProject)
	})

	api.Route("/clients", func(clients fiber.Router) {
		clients.Get("/", controller.GetClients)
		clients.Get("/:clientId", controller.GetClient)
		clients.Post("/", controller.CreateClient)
		clients.Put("/:clientId", controller.UpdateClient)
		clients.Delete("/:clientId", controller.DeleteClient)
	})
}
