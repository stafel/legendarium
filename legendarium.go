package main

import (
	"log"
	"os"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type LegendCharacter struct {
	gorm.Model

	CharacterSheet string
	Name           string
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

func main() {
	db, err := gorm.Open(sqlite.Open("legends.db"), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}

	db.AutoMigrate(&LegendCharacter{})

	cs, err := readCharsheet("../Far'dei.json")
	if err != nil {
		log.Fatal(err)
	}

	newLegendChar := &LegendCharacter{
		CharacterSheet: cs,
	}

	db.Create(newLegendChar)
}
