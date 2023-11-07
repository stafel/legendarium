package provider

import (
	"testing"

	"github.com/stafel/legendarium/model"
)

func TestGetMilestones(t *testing.T) {
	var testedCharId uint = 1

	s := &SqliteProvider{}
	s.Connect()

	ms, err := s.GetMilestones(&model.LegendMilestone{LegendCharacterID: testedCharId})
	if err != nil {
		t.Error(err)
	}

	for _, milestone := range ms {
		if milestone.LegendCharacterID != testedCharId {
			t.Error("Wrong character reference id found")
		}
	}
}
