package main

import (
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type LegendCharacter struct {
	gorm.Model

	Name        string `gorm:"uniqueIndex"`
	Description string
	Milestones  []LegendMilestone
}

// Contains Character sheet and metadata valid in a point of time
type LegendMilestone struct {
	gorm.Model
	LegendCharacterID uint // foreign key for the has many reference

	ReferenceDate  time.Time
	CharacterSheet string
	Description    string
}

// Read full file content as string
func readCharsheet(path string) (string, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}

	return string(content), nil
}

// extracted info by getPaths
type PathInfo struct {
	Name    string
	Path    string
	RefDate time.Time
	Note    string
}

func getPaths(basepath string) []PathInfo {
	dirs, err := os.ReadDir(basepath)
	if err != nil {
		log.Fatal(err)
	}

	paths := make([]PathInfo, 0)

	for _, dir := range dirs {
		if !dir.IsDir() {
			if !strings.Contains(dir.Name(), "_") || !strings.Contains(dir.Name(), ".json") {
				continue
			}

			fullPath := filepath.Join(basepath, dir.Name())

			substrings := strings.Split(strings.Split(dir.Name(), ".")[0], "_")

			dateOnlyString := substrings[1][4:] + "-" + substrings[1][2:4] + "-" + substrings[1][0:2] // "2006-01-02"

			refDate, err := time.Parse(time.DateOnly, dateOnlyString)
			if err != nil {
				log.Fatal(err)
			}

			notes := strings.Join(substrings[2:], " ")

			paths = append(paths, PathInfo{
				Name:    substrings[0],
				Path:    fullPath,
				RefDate: refDate,
				Note:    notes,
			})
		}
	}

	return paths
}

// returns a list paths of all jsons in the folder
func loadFolder(folderPath string) []*LegendCharacter {

	legendChars := make([]*LegendCharacter, 0)

	for i, pInfo := range getPaths(folderPath) {
		log.Printf("%v %v", i, pInfo)

		var currentLegendChar *LegendCharacter
		currentLegendChar = nil
		for _, lc := range legendChars {
			if lc.Name == pInfo.Name {
				currentLegendChar = lc
			}
		}
		if currentLegendChar == nil {
			currentLegendChar = &LegendCharacter{
				Name: pInfo.Name,
			}
			legendChars = append(legendChars, currentLegendChar)
		}

		cs, err := readCharsheet(pInfo.Path)
		if err != nil {
			log.Fatal(err)
		}

		newMilestone := &LegendMilestone{
			CharacterSheet: cs,
			ReferenceDate:  pInfo.RefDate,
			Description:    pInfo.Note,
		}

		currentLegendChar.Milestones = append(currentLegendChar.Milestones, *newMilestone)
	}

	return legendChars
}

// loads all legends from the folder as jsons into the db. sorts out duplicates by character name and milestone json data
func loadFolderIntoDb(folderPath string, db *gorm.DB) {
	for _, legend := range loadFolder(folderPath) {

		var existingChar LegendCharacter
		err := db.Where(&LegendCharacter{Name: legend.Name}).Preload("Milestones").First(&existingChar).Error
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

		db.Save(legend)
	}
}

func main() {
	db, err := gorm.Open(sqlite.Open("legends.db"), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}

	db.AutoMigrate(&LegendCharacter{})
	db.AutoMigrate(&LegendMilestone{})

	loadFolderIntoDb("/data/MyDocumentsRPG/earthdawn", db)
}
