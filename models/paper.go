package models

import "gorm.io/gorm"

type Paper struct {
	gorm.Model
	Title     string   `json:"title" gorm:"size:80;unique_index;not null" form:"title"`
	Abstract  string   `json:"abstract" gorm:"not null" form:"abstract"`
	Content   []byte   `gorm:"not null" form:"content"`
	Authors   []Author `gorm:"many2many:paper_authors" json:"author" form:"author"`
	FacultyID uint     `gorm:"not null" json:"faculty_id" form:"faculty_id"`
}
