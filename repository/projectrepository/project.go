package projectrepository

import (
	"datcha/repository/repositorycommon"
	"strconv"

	"gorm.io/gorm"
)

type RepProject struct {
	ID        uint
	Name      string `gorm:"uniqueIndex:project_name_idx"`
	ImagePath string
	OwnerID   uint `gorm:"uniqueIndex:project_name_idx"`
}

func NewRepProject(name string, ownerId uint) *RepProject {
	project := RepProject{
		Name:    name,
		OwnerID: ownerId,
	}
	return &project
}

func (project *RepProject) AfterUpdate(tx *gorm.DB) (err error) {
	repositorycommon.SendNotify(tx, "project updated : "+project.Name)
	return nil
}

func (project *RepProject) AfterCreate(tx *gorm.DB) error {
	repositorycommon.SendNotify(tx, "project added "+strconv.FormatUint(uint64(project.ID), 10))
	return nil
}

func (project *RepProject) AfterDelete(tx *gorm.DB) error {
	repositorycommon.SendNotify(tx, "project deleted "+strconv.FormatUint(uint64(project.ID), 10))
	return nil
}
