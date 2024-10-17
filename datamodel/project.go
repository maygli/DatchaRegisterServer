package datamodel

import (
	"datcha/servercommon"
	"strconv"

	"gorm.io/gorm"
)

type Project struct {
	ID               uint
	Name             string `gorm:"uniqueIndex:project_name_idx"`
	ImagePath        string
	OwnerID          uint `gorm:"uniqueIndex:project_name_idx"`
	Owner            *User
	Devices          []Device
	ProjectCardItems []ProjectCardItem
}

func NewProject(name string, owner *User) *Project {
	project := Project{
		Name:  name,
		Owner: owner,
	}
	return &project
}

func (project *Project) AfterUpdate(tx *gorm.DB) (err error) {
	servercommon.SendNotify(tx, "project updated : "+project.Name)
	return nil
}

func (project *Project) AfterCreate(tx *gorm.DB) error {
	servercommon.SendNotify(tx, "project added "+strconv.FormatUint(uint64(project.ID), 10))
	return nil
}

func (project *Project) AfterDelete(tx *gorm.DB) error {
	servercommon.SendNotify(tx, "project deleted "+strconv.FormatUint(uint64(project.ID), 10))
	return nil
}
