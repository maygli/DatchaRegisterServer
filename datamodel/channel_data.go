package datamodel

import (
	"datcha/servercommon"
	"time"

	"gorm.io/gorm"
)

type ChannelData struct {
	ID        uint `gorm:"primarykey"`
	CreatedAt time.Time
	Data      string
	Unit      string
	ChannelID uint
	Channel   *Channel
}

func (channelData *ChannelData) AfterUpdate(tx *gorm.DB) (err error) {
	servercommon.SendNotify(tx, "channel data changed")
	return nil
}

func (channelData *ChannelData) AfterCreate(tx *gorm.DB) error {
	servercommon.SendNotify(tx, "channel data created")
	return nil
}

func (channelData *ChannelData) AfterDelete(tx *gorm.DB) error {
	servercommon.SendNotify(tx, "channel data removed")
	return nil
}
