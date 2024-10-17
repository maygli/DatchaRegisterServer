package datamodel

import "time"

type Command struct {
	ID        uint `gorm:"primarykey"`
	CreatedAt time.Time
	Key       string `gorm:"type:varchar(100);uniqueIndex:command_key_idx"`
	Value     string
	DeviceID  uint
	Device    *Device
}
