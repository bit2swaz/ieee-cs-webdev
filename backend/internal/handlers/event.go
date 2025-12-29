package handlers

import (
	"time"

	"github.com/bit2swaz/ieee-cs-webdev/backend/internal/database"
	"github.com/bit2swaz/ieee-cs-webdev/backend/internal/models"
	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
)

type CreateEventInput struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Location    string `json:"location"`
	MaxCapacity int    `json:"max_capacity"`
	Date        string `json:"date"` // Format: YYYY-MM-DDTHH:MM:SSZ
}

func CreateEvent(c fiber.Ctx) error {
	// 1. Parse Input
	var input CreateEventInput
	if err := c.Bind().Body(&input); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid Input"})
	}

	// 2. Get User's Org ID from Middleware (The Magic)
	orgIDStr := c.Locals("org_id").(string)
	orgID, _ := uuid.Parse(orgIDStr)

	// 3. Parse Date
	parsedDate, err := time.Parse(time.RFC3339, input.Date)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid Date Format. Use ISO8601"})
	}

	// 4. Create Event Scoped to Org
	event := models.Event{
		Title:          input.Title,
		Description:    input.Description,
		Location:       input.Location,
		MaxCapacity:    input.MaxCapacity,
		Date:           parsedDate,
		OrganizationID: orgID, // <--- Secured Multi-Tenancy
	}

	if err := database.DB.Create(&event).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Could not create event"})
	}

	return c.JSON(event)
}

func GetEvents(c fiber.Ctx) error {
	// 1. Get User's Org ID
	orgIDStr := c.Locals("org_id").(string)

	var events []models.Event

	// 2. Fetch Only Events belonging to this Org
	// This ensures Org A never sees Org B's events
	if err := database.DB.Where("organization_id = ?", orgIDStr).Find(&events).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Could not fetch events"})
	}

	return c.JSON(events)
}
