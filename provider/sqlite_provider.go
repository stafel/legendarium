package provider

import (
	"log"

	"github.com/stafel/legendarium/model"
	"github.com/stafel/legendarium/other"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type SqliteProvider struct {
	db *gorm.DB
}

func (s *SqliteProvider) Connect() {
	db, err := gorm.Open(sqlite.Open("legends.db"), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}

	s.db = db

	db.AutoMigrate(&model.LegendCharacter{})
	db.AutoMigrate(&model.LegendMilestone{})
}

// loads all legends from the folder as jsons into the db. sorts out duplicates by character name and milestone json data
func (s *SqliteProvider) MigrateFromFolder(folderPath string) {
	// loads all legends from the folder as jsons into the db. sorts out duplicates by character name and milestone json data
	for _, legend := range other.LoadFolder(folderPath) {

		var existingChar model.LegendCharacter
		err := s.db.Where(&model.LegendCharacter{Name: legend.Name}).Preload("Milestones").First(&existingChar).Error
		if err == nil {
			log.Printf("Char %v already exists with ID %v", legend.Name, existingChar.ID)
			legend.ID = existingChar.ID
			legend.CreatedAt = existingChar.CreatedAt
			legend.DeletedAt = existingChar.DeletedAt

			for i, newMstone := range legend.Milestones {
				for _, existingMstone := range existingChar.Milestones {
					if newMstone.CharacterSheet == existingMstone.CharacterSheet {
						log.Printf("Char %v Milestone %v already exists with ID %v", legend.Name, i, existingMstone.ID)
						legend.Milestones[i].ID = existingMstone.ID
						legend.Milestones[i].CreatedAt = existingMstone.CreatedAt
						legend.Milestones[i].DeletedAt = existingMstone.DeletedAt
					}
				}
			}
		}

		s.db.Save(legend)
	}
}

func (s *SqliteProvider) GetCharacters() []model.LegendCharacter {
	var characters []model.LegendCharacter
	s.db.Model(&model.LegendCharacter{}).Find(&characters)
	return characters
}

func (s *SqliteProvider) GetCharacter(searchReferenceChar *model.LegendCharacter) (model.LegendCharacter, error) {
	var existingChar model.LegendCharacter
	err := s.db.Where(searchReferenceChar).First(&existingChar).Error
	if err != nil {
		return existingChar, err
	}
	return existingChar, nil
}
