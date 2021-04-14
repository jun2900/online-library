package models

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Email    string `json:"email" gorm:"size:45;unique_index;not null"`
	Password string `gorm:"size:80;not null" json:"password"`
}
