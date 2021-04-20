package models

import "gorm.io/gorm"

type Paper struct {
	gorm.Model
	Title     string   `json:"title" gorm:"size:80;unique_index;not null" validate:"required"`
	Abstract  string   `json:"abstract" gorm:"not null" validate:"required"`
	Content   []byte   `json:"content" gorm:"not null" validate:"required"`
	Authors   []Author `json:"authors" gorm:"many2many:paper_authors" validate:"required"`
	FacultyID uint     `json:"facultyId" gorm:"not null" validate:"required"`
}
