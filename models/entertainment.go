package models

import "gorm.io/gorm"

type Entertainment struct {
	gorm.Model
	Name string `gorm:"type:varchar(60); not null; unique" json:"name"`
	ImgUrl	string	`gorm:"type:varchar(512); unique;" json:"imgUrl"`
	//Artists []Artist
}
