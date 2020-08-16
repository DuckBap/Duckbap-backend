package models

import "gorm.io/gorm"

type Artist struct {
	gorm.Model
	EntertainmentID uint `gorm:"not null" json:"entertainmentId"`
	Name string `gorm:"type:varchar(60); not null" json:"name"`
	ImgUrl string `gorm:"type:varchar(512); unique; not null" json:"imgUrl"`
	Users []User `gorm:"many2many:bookmarks"`
	Entertainment Entertainment
}
