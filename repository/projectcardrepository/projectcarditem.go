package projectcardrepository

import "time"

type RepProjectCardItem struct {
	ID        uint `gorm:"primarykey"`
	CreatedAt time.Time
	ChannelID uint `gorm:"foreignKey:ProjectCardChannelRefer"`
	ProjectID uint `gorm:"foreignKey:ProjectCardProjectRefer"`
}
