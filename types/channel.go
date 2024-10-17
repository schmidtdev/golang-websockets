package types

import "gorm.io/gorm"

type Channel struct {
	gorm.Model
	ID   uint   `json:"id"`
	Name string `json:"name"`
}
