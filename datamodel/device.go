package datamodel

type Device struct {
	ID        uint   `gorm:"primarykey"`
	Name      string `gorm:"type:varchar(100);uniqueIndex:device_name_idx"`
	ProjectID uint   `gorm:"uniqueIndex:device_name_idx"`
	Project   *Project
}
