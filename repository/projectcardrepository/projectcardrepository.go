package projectcardrepository

import (
	"datcha/datamodel"
	"datcha/servercommon"

	"gorm.io/gorm"
)

type ProjectCardRepository struct {
	database *gorm.DB
	notifier servercommon.INotifier
}

func NewProjectCardRepository(db *gorm.DB, notifier servercommon.INotifier) (ProjectCardRepositorier, error) {
	rep := ProjectCardRepository{
		database: db,
		notifier: notifier,
	}
	err := rep.database.AutoMigrate(&RepProjectCardItem{})
	return &rep, err
}

func (repo *ProjectCardRepository) AddProjectCardItem(item *datamodel.ProjectCardItem) error {
	repItem, err := getNewRepProjectCardItem(item)
	if err != nil {
		return err
	}
	result := repo.database.Set(servercommon.NOTIFIER_KEY, repo.notifier).Create(repItem)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (repo *ProjectCardRepository) RemoveProjectCardItem(projectCardItemId uint) error {
	result := repo.database.Set(servercommon.NOTIFIER_KEY, repo.notifier).Delete(&RepProjectCardItem{}, projectCardItemId)
	return result.Error
}

func (repo *ProjectCardRepository) GetProjectCardItems(projectId uint) ([]*datamodel.ProjectCardItem, error) {
	cardItems := []RepProjectCardItem{}
	result := repo.database.Where("project_id=?", projectId).Find(&cardItems)
	if result.Error != nil {
		return []*datamodel.ProjectCardItem{}, result.Error
	}
	return getNewDMProjectCardItems(cardItems)
}
