package model

import (
	"time"

	"gorm.io/gorm"
)

// Contains Character sheet and metadata valid in a point of time
type LegendMilestone struct {
	gorm.Model
	LegendCharacterID uint // foreign key for the has many reference

	ReferenceDate  time.Time
	CharacterSheet string
	Description    string
}
