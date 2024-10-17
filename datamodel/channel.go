package datamodel

type Channel struct {
	ID       uint   `gorm:"primarykey"`
	Name     string `gorm:"type:varchar(100);uniqueIndex:channel_name_idx"`
	DeviceID uint   `gorm:"uniqueIndex:channel_name_idx"`
	Device   *Device
}
