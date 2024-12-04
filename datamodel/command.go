package datamodel

import "time"

type Command struct {
	ID        uint
	CreatedAt time.Time
	Key       string
	Value     string
	DeviceID  uint
}
