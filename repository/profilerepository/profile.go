package profilerepository

import (
	"datcha/repository/repositorycommon"

	"gorm.io/gorm"
)

type RepProfile struct {
	ID     uint   `gorm:"primarykey"`
	Name   string `gorm:"type:varchar(100)"`
	Avatar string `gorm:"type:varchar(2048)"`
	UserID uint   `gorm:"uniqueIndex:profile_user_idx"`
}

func (profile *RepProfile) AfterUpdate(tx *gorm.DB) (err error) {
	repositorycommon.SendNotify(tx, "profile data changed")
	return nil
}

func (profile *RepProfile) AfterCreate(tx *gorm.DB) error {
	repositorycommon.SendNotify(tx, "channel data changed")
	return nil
}

func (profile *RepProfile) AfterDelete(tx *gorm.DB) error {
	repositorycommon.SendNotify(tx, "channel data changed")
	return nil
}
