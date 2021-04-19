package models

import "gorm.io/gorm"

type Paper struct {
	gorm.Model
	Title     string   `json:"title" gorm:"size:80;unique_index;not null"`
	Abstract  string   `json:"abstract" gorm:"not null"`
	Content   []byte   `json:"content" gorm:"not null"`
	Authors   []Author `json:"authors" gorm:"many2many:paper_authors"`
	FacultyID uint     `json:"facultyId" gorm:"not null"`
}
