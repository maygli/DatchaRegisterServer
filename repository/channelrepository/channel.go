package channelrepository

import (
	"datcha/repository/repositorycommon"
	"time"

	"gorm.io/gorm"
)

type RepChannel struct {
	ID       uint   `gorm:"primarykey"`
	Name     string `gorm:"type:varchar(100);uniqueIndex:channel_name_idx"`
	DeviceID uint   `gorm:"uniqueIndex:channel_name_idx"`
}

type RepChannelData struct {
	ID        uint `gorm:"primarykey"`
	CreatedAt time.Time
	Data      string
	Unit      string
	ChannelID uint
}

func (channel *RepChannel) AfterUpdate(tx *gorm.DB) (err error) {
	repositorycommon.SendNotify(tx, "channel data changed")
	return nil
}

func (channel *RepChannel) AfterCreate(tx *gorm.DB) error {
	repositorycommon.SendNotify(tx, "channel data changed")
	return nil
}

func (channel *RepChannel) AfterDelete(tx *gorm.DB) error {
	repositorycommon.SendNotify(tx, "channel data changed")
	return nil
}

func (channelData *RepChannelData) AfterUpdate(tx *gorm.DB) (err error) {
	repositorycommon.SendNotify(tx, "channel data changed")
	return nil
}

func (channelData *RepChannelData) AfterCreate(tx *gorm.DB) error {
	repositorycommon.SendNotify(tx, "channel data changed")
	return nil
}

func (channelData *RepChannelData) AfterDelete(tx *gorm.DB) error {
	repositorycommon.SendNotify(tx, "channel data changed")
	return nil
}
