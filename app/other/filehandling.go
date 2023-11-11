package other

import (
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/stafel/legendarium/model"
)

// extracted info by getPaths
type PathInfo struct {
	Name    string
	Path    string
	RefDate time.Time
	Note    string
}

// get all json file paths from a folder
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

// Read full file content as string
func readCharsheet(path string) (string, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}

	return string(content), nil
}

// returns a list paths of all jsons in the folder
func LoadFolder(folderPath string) []*model.LegendCharacter {

	legendChars := make([]*model.LegendCharacter, 0)

	for i, pInfo := range getPaths(folderPath) {
		log.Printf("%v %v", i, pInfo)

		var currentLegendChar *model.LegendCharacter
		currentLegendChar = nil
		for _, lc := range legendChars {
			if lc.Name == pInfo.Name {
				currentLegendChar = lc
			}
		}
		if currentLegendChar == nil {
			currentLegendChar = &model.LegendCharacter{
				Name: pInfo.Name,
			}
			legendChars = append(legendChars, currentLegendChar)
		}

		cs, err := readCharsheet(pInfo.Path)
		if err != nil {
			log.Fatal(err)
		}

		newMilestone := &model.LegendMilestone{
			CharacterSheet: cs,
			ReferenceDate:  pInfo.RefDate,
			Description:    pInfo.Note,
		}

		currentLegendChar.Milestones = append(currentLegendChar.Milestones, *newMilestone)
	}

	return legendChars
}
