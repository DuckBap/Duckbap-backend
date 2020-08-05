package models

import "gorm.io/gorm"

type Artist struct {
	gorm.Model
	EntertainmentID uint `gorm:"not null"`
	Name string `gorm:"varchar(60); not null"`
	ImgUrl string `gorm:"varchar(255); unique; not null"`
	Users []User `gorm:"many2many:bookmarks"`
	Entertainment Entertainment
}
