package projectrepository

import (
	"datcha/datamodel"
	"datcha/servercommon"
	"errors"

	"gorm.io/gorm"
)

type ProjectRepository struct {
	database *gorm.DB
	notifier servercommon.INotifier
}

func NewProjectRepository(db *gorm.DB, notifier servercommon.INotifier) (ProjectRepositorier, error) {
	rep := ProjectRepository{
		database: db,
		notifier: notifier,
	}
	err := rep.database.AutoMigrate(&RepProject{})
	return &rep, err
}

func (repo *ProjectRepository) AddProject(project *datamodel.Project) error {
	repproject, err := getNewRepProject(project)
	if err != nil {
		return err
	}
	result := repo.database.Set(servercommon.NOTIFIER_KEY, repo.notifier).Create(repproject)
	if result.Error != nil {
		return errors.New(servercommon.ERROR_INTERNAL)
	}
	err = updateDMProject(repproject, project)
	return err
}

func (repo *ProjectRepository) FindProject(projectId uint) (*datamodel.Project, error) {
	project := RepProject{}
	result := repo.database.First(&project, projectId)
	if result.Error != nil {
		return &datamodel.Project{}, result.Error
	}
	dmProject, err := getNewDMProject(&project)
	return dmProject, err
}

func (repo *ProjectRepository) UpdateProject(project *datamodel.Project) error {
	repProject, err := getNewRepProject(project)
	if err != nil {
		return err
	}
	result := repo.database.Set(servercommon.NOTIFIER_KEY, repo.notifier).Save(repProject)
	return result.Error
}

func (repo *ProjectRepository) GetOwnerProjects(userId uint) ([]*datamodel.Project, error) {
	projects := []RepProject{}
	result := repo.database.Where("owner_id = ?", userId).Find(&projects)
	if result.Error != nil {
		return []*datamodel.Project{}, result.Error
	}
	dmProjects, err := getNewDMProjectsList(projects)
	return dmProjects, err
}

func (repo *ProjectRepository) DeleteProject(projectId uint) error {
	result := repo.database.Set(servercommon.NOTIFIER_KEY, repo.notifier).Delete(&RepProject{}, projectId)
	return result.Error
}
