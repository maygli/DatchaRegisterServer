package datamodel

import (
	"time"
)

type ChannelData struct {
	ID        uint
	CreatedAt time.Time
	Data      string
	Unit      string
	ChannelID uint
}
