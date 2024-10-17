package types

import "gorm.io/gorm"

type Subscription struct {
	gorm.Model
	ID        uint    `json:"id"`
	IP        string  `json:"ip"`
	ChannelID uint    `json:"channel_id"`
	Channel   Channel `json:"channel"`
}
