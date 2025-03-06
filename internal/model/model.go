package model

import "gorm.io/gorm"

type Post struct {
	gorm.Model
	Title    string
	Content  string    `gorm:"type:text"`
	Comments []Comment `gorm:"foreignKey:PostID"`
}

type Comment struct {
	gorm.Model
	Content string `gorm:"type:text"`
	PostID  uint
}
