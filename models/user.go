package models

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Email    string `json:"email" gorm:"size:45"`
	Password string `gorm:"size:80"`
}
