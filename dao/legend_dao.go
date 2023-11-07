package dao

import (
	"github.com/stafel/legendarium/model"
	"github.com/stafel/legendarium/provider"
)

type LegendDao struct {
	provider *provider.SqliteProvider
}

func DefaultConnect() *LegendDao {
	l := &LegendDao{}
	prov := &provider.SqliteProvider{}
	prov.Connect()
	l.provider = prov
	return l
}

func (l *LegendDao) Connect(prov *provider.SqliteProvider) {
	l.provider = prov
}

func (l *LegendDao) MigrateFromFolder(folderPath string) {
	l.provider.MigrateFromFolder(folderPath)
}

func (l *LegendDao) GetCharacters() []model.LegendCharacter {
	return l.provider.GetCharacters()
}

func (l *LegendDao) GetCharacterById(id uint) (model.LegendCharacter, error) {
	search := &model.LegendCharacter{}
	search.ID = id
	return l.provider.GetCharacter(search)
}

func (l *LegendDao) GetMilestonesForCharacterId(id uint) ([]model.LegendMilestone, error) {
	return l.provider.GetMilestones(&model.LegendMilestone{LegendCharacterID: id})
}

func (l *LegendDao) GetLatestMilestonesForCharacterId(id uint) (model.LegendMilestone, error) {
	return l.provider.GetLatestMilestone(&model.LegendMilestone{LegendCharacterID: id})
}

func (l *LegendDao) GetMilestoneForMilestoneId(id uint) (model.LegendMilestone, error) {
	referenceMilestone := &model.LegendMilestone{}
	referenceMilestone.ID = id

	ms, err := l.provider.GetMilestones(referenceMilestone)
	if err != nil {
		return *referenceMilestone, err
	}

	return ms[0], nil
}
