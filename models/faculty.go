package models

import "gorm.io/gorm"

type Faculty struct {
	gorm.Model
	Name string `json:"name" gorm:"size:15"`
}
