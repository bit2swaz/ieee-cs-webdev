package handlers

import (
	"github.com/bit2swaz/ieee-cs-webdev/backend/internal/database"
	"github.com/bit2swaz/ieee-cs-webdev/backend/internal/models"
	"github.com/bit2swaz/ieee-cs-webdev/backend/internal/utils"
	"github.com/gofiber/fiber/v3"
)

// SignupInput defines what we expect from the frontend
type SignupInput struct {
	OrgName  string `json:"org_name"`
	Domain   string `json:"domain"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginInput struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func Register(c fiber.Ctx) error {
	var input SignupInput
	if err := c.Bind().Body(&input); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid Input"})
	}

	// 1. Create Organization
	org := models.Organization{
		Name:   input.OrgName,
		Domain: input.Domain,
	}
	// Transaction ensures both Org and User are created, or neither
	tx := database.DB.Begin()

	if err := tx.Create(&org).Error; err != nil {
		tx.Rollback()
		return c.Status(500).JSON(fiber.Map{"error": "Could not create Organization (Domain might be taken)"})
	}

	// 2. Hash Password
	hash, err := utils.HashPassword(input.Password)
	if err != nil {
		tx.Rollback()
		return c.Status(500).JSON(fiber.Map{"error": "Password hashing failed"})
	}

	// 3. Create User (Admin of that Org)
	user := models.User{
		Name:           input.Name,
		Email:          input.Email,
		Password:       hash,
		Role:           "admin", // First user is always admin
		OrganizationID: org.ID,
	}

	if err := tx.Create(&user).Error; err != nil {
		tx.Rollback()
		return c.Status(500).JSON(fiber.Map{"error": "Could not create User"})
	}

	tx.Commit()

	return c.JSON(fiber.Map{"message": "Registration successful", "org_id": org.ID})
}

func Login(c fiber.Ctx) error {
	var input LoginInput
	if err := c.Bind().Body(&input); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid Input"})
	}

	var user models.User
	// Find user by email
	if err := database.DB.Where("email = ?", input.Email).First(&user).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "User not found"})
	}

	// Check password
	if !utils.CheckPasswordHash(input.Password, user.Password) {
		return c.Status(401).JSON(fiber.Map{"error": "Invalid Password"})
	}

	// Generate JWT
	token, err := utils.GenerateToken(user.ID, user.OrganizationID, user.Role)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Could not login"})
	}

	return c.JSON(fiber.Map{"token": token})
}
