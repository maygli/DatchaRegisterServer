package repository

import (
	"datcha/datamodel"
	"datcha/servercommon"
	"errors"
)

func (repo *DatchaSqlRepository) AddProjectCardItem(item *datamodel.ProjectCardItem) error {
	result := repo.database.Set(servercommon.NOTIFIER_KEY, repo.notifier).Create(item)
	if result.Error != nil {
		return errors.New(servercommon.ERROR_INTERNAL)
	}
	return nil
}

func (repo *DatchaSqlRepository) RemoveProjectCardItem(projectCardItemId uint) error {
	repo.database.Set(servercommon.NOTIFIER_KEY, repo.notifier).Delete(&datamodel.ProjectCardItem{}, projectCardItemId)
	return nil
}

func (repo *DatchaSqlRepository) GetProjectCardItems(projectId uint) ([]datamodel.ProjectCardItem, error) {
	cardItems := []datamodel.ProjectCardItem{}
	result := repo.database.Where("project_id=?", projectId).Find(&cardItems)
	return cardItems, result.Error
}
