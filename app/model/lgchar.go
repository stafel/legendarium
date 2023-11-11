package model

import "gorm.io/gorm"

type LegendCharacter struct {
	gorm.Model

	Name        string `gorm:"uniqueIndex"`
	Description string
	Milestones  []LegendMilestone
}
