package projectrepository

import "datcha/datamodel"

type ProjectRepositorier interface {
	AddProject(project *datamodel.Project) error
	FindProject(projectId uint) (*datamodel.Project, error)
	UpdateProject(project *datamodel.Project) error
	DeleteProject(projectId uint) error
	GetOwnerProjects(userId uint) ([]*datamodel.Project, error)
}
