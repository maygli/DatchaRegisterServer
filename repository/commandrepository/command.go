package commandrepository

import "time"

type RepCommand struct {
	ID        uint `gorm:"primarykey"`
	CreatedAt time.Time
	Key       string `gorm:"type:varchar(100);uniqueIndex:command_key_idx"`
	Value     string
	DeviceID  uint `gorm:"uniqueIndex:command_key_idx"`
}
