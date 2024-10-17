package repository

import (
	"datcha/datamodel"
	"datcha/servercommon"
	"errors"
)

func (repo *DatchaSqlRepository) AddProject(project *datamodel.Project) error {
	result := repo.database.Set(servercommon.NOTIFIER_KEY, repo.notifier).Create(project)
	if result.Error != nil {
		return errors.New(servercommon.ERROR_INTERNAL)
	}
	return nil
}

func (repo *DatchaSqlRepository) RemoveProject(projectId uint) error {
	return nil
}

func (repo *DatchaSqlRepository) FindProject(projectId uint) (datamodel.Project, error) {
	project := datamodel.Project{}
	result := repo.database.Preload("Devices").Preload("ProjectCardItems").First(&project, projectId)
	return project, result.Error
}

func (repo *DatchaSqlRepository) UpdateProject(project *datamodel.Project) error {
	repo.database.Set(servercommon.NOTIFIER_KEY, repo.notifier).Save(project)
	return nil
}

func (repo *DatchaSqlRepository) GetOwnerProjects(userId uint) ([]datamodel.Project, error) {
	projects := []datamodel.Project{}
	result := repo.database.Where("owner_id = ?", userId).Preload("Devices").Preload("ProjectCardItems").Find(&projects)
	return projects, result.Error
}

func (repo *DatchaSqlRepository) DeleteProject(projectId uint) error {
	result := repo.database.Set(servercommon.NOTIFIER_KEY, repo.notifier).Delete(&datamodel.Project{}, projectId)
	return result.Error
}
