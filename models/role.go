package models

import "gorm.io/gorm"

type Role struct {
	gorm.Model
	Name string `json:"name" gorm:"size:45;unique_index;not null" validate:"required"`
	User []User
}
