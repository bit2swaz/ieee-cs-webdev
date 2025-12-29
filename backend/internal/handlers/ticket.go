package handlers

import (
	"fmt"

	"github.com/bit2swaz/ieee-cs-webdev/backend/internal/database"
	"github.com/bit2swaz/ieee-cs-webdev/backend/internal/models"
	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func BookTicket(c fiber.Ctx) error {
	// 1. Get Event ID from params
	eventIDStr := c.Params("id")
	eventID, err := uuid.Parse(eventIDStr)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid Event ID"})
	}

	// 2. Get User ID from Middleware
	userIDStr := c.Locals("user_id").(string)
	userID, _ := uuid.Parse(userIDStr)

	// 3. START TRANSACTION (The Critical Part)
	err = database.DB.Transaction(func(tx *gorm.DB) error {
		var event models.Event

		// 4. LOCK THE ROW
		// "clause.Locking{Strength: 'UPDATE'}" translates to "SELECT * FROM events FOR UPDATE"
		// This tells Postgres: "Don't let anyone else read/write this specific row until I'm done."
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).First(&event, "id = ?", eventID).Error; err != nil {
			return err // Event not found
		}

		// 5. Check if user already booked (Optional, but good practice)
		var existingTicket models.Ticket
		if err := tx.Where("event_id = ? AND user_id = ?", eventID, userID).First(&existingTicket).Error; err == nil {
			return fmt.Errorf("already booked")
		}

		// 6. Check Logic (Inside the lock)
		if event.TicketsSold >= event.MaxCapacity {
			return fmt.Errorf("sold out")
		}

		// 7. Create Ticket
		ticket := models.Ticket{
			EventID:    eventID,
			UserID:     userID,
			TicketCode: uuid.New().String(), // Unique code
			Status:     "booked",
		}

		if err := tx.Create(&ticket).Error; err != nil {
			return err
		}

		// 8. Update Event Counter
		event.TicketsSold++
		if err := tx.Save(&event).Error; err != nil {
			return err
		}

		// Transaction commits automatically if we return nil
		return nil
	})

	// Handle Transaction Errors
	if err != nil {
		if err.Error() == "sold out" {
			return c.Status(409).JSON(fiber.Map{"error": "Event is sold out"})
		}
		if err.Error() == "already booked" {
			return c.Status(409).JSON(fiber.Map{"error": "You already have a ticket"})
		}
		return c.Status(500).JSON(fiber.Map{"error": "Booking failed"})
	}

	return c.Status(201).JSON(fiber.Map{"message": "Ticket booked successfully!"})
}
