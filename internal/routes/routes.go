package routes

import (
	"fold/internal/controllers"
	"fold/internal/middlewares"

	"github.com/gofiber/fiber/v2"
)

func SetupRoutes() *fiber.App {
	app := fiber.New()

	//Middleware
	app.Use(middlewares.SetRequestID())
	app.Use(middlewares.CORS())

	//Routes (unauthenticated)
	api := app.Group("/api/v1")
	api.Get("/projectForUser/:userName", controllers.GetProjectsForUser)
	api.Get("/projectAndUserForHashtag/:hashtag", controllers.ProjectAndUserForHashtag)
	api.Get("/FuzzySearchProject/:tag", controllers.FuzzySearchProject)

	return app
}
