package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Base model for ID and Timestamps
type Base struct {
	ID        uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primary_key"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

// 1. Organization (The Tenant)
// Every user and event belongs to one of these.
type Organization struct {
	Base
	Name   string `gorm:"not null"`
	Domain string `gorm:"unique;not null"` // e.g., "ieee.sastra.edu"
	ApiKey string `gorm:"unique"`          // For API access (Level 3)
}

// 2. User
// Scoped to an Organization.
type User struct {
	Base
	Name           string `gorm:"not null"`
	Email          string `gorm:"not null"` // Unique per Org, handled by logic or composite index
	Password       string `gorm:"not null"`
	Role           string `gorm:"default:'user'"` // 'admin', 'user'
	OrganizationID uuid.UUID
	Organization   Organization `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

// 3. Event
// The main entity. Can be a standalone event or a "Fest" containing sub-events.
type Event struct {
	Base
	Title          string `gorm:"not null"`
	Description    string
	Date           time.Time
	Location       string
	MaxCapacity    int  `gorm:"not null"`
	TicketsSold    int  `gorm:"default:0"`
	IsFest         bool `gorm:"default:false"` // Level 2: Fest style
	OrganizationID uuid.UUID
	Organization   Organization
}

// 4. SubEvent (Level 2)
// Example: "Hackathon" inside the "Tech Fest".
type SubEvent struct {
	Base
	Title     string `gorm:"not null"`
	StartTime time.Time
	EndTime   time.Time
	EventID   uuid.UUID
	Event     Event
}

// 5. Ticket
// Links a User to an Event.
type Ticket struct {
	Base
	TicketCode string `gorm:"unique;not null"`
	Status     string `gorm:"default:'booked'"` // booked, cancelled, checked-in
	EventID    uuid.UUID
	Event      Event
	UserID     uuid.UUID
	User       User
}
