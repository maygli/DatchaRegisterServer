package projectcardrepository

import "datcha/datamodel"

type ProjectCardRepositorier interface {
	AddProjectCardItem(cardItem *datamodel.ProjectCardItem) error
	RemoveProjectCardItem(projectCardItemId uint) error
	GetProjectCardItems(projectId uint) ([]*datamodel.ProjectCardItem, error)
}
