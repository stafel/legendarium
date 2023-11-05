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

func (l *LegendDao) GetCharacterById(id uint) model.LegendCharacter {
	search := &model.LegendCharacter{}
	search.ID = id
	c, _ := l.provider.GetCharacter(search)
	return c
}
