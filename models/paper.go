package models

import "gorm.io/gorm"

type Paper struct {
	gorm.Model
	Title    string `json:"title" gorm:"size:80"`
	Abstract string
	Content  []byte
	Authors  []Author `gorm:"many2many:paper_authors"`
}
