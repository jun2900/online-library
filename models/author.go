package models

import "gorm.io/gorm"

type Author struct {
	gorm.Model
	Name   string `json:"name" gorm:"size:50"`
	Papers []Paper
}
