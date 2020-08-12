package models

import "gorm.io/gorm"

type Entertainment struct {
	gorm.Model
	Name string `gorm:"varchar(60); not null; unique" json:"name"`
	ImgUrl	string	`gorm:"varchar(255); unique;" json:"imgUrl"`
	//Artists []Artist
}
