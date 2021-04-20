package models

import "gorm.io/gorm"

type Author struct {
	gorm.Model
	Name string `json:"name" gorm:"unique;not null;size:50"`
}
