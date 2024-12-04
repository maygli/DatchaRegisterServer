package datamodel

import "time"

type ProjectCardItem struct {
	ID        uint
	CreatedAt time.Time
	ChannelID uint
	ProjectID uint
}
