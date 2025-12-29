package routes

import (
	"github.com/bit2swaz/ieee-cs-webdev/backend/internal/handlers"
	"github.com/bit2swaz/ieee-cs-webdev/backend/internal/middleware"
	"github.com/gofiber/fiber/v3"
)

func SetupRoutes(app *fiber.App) {
	api := app.Group("/api")

	// Public Routes
	auth := api.Group("/auth")
	auth.Post("/register", handlers.Register)
	auth.Post("/login", handlers.Login)

	// Protected Routes (Require JWT)
	event := api.Group("/events")
	event.Use(middleware.Protected()) // Lock the door

	event.Post("/", handlers.CreateEvent)
	event.Get("/", handlers.GetEvents)

	event.Post("/:id/book", handlers.BookTicket)
}
