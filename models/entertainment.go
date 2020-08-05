package models

import "gorm.io/gorm"

type Entertainment struct {
	gorm.Model
	Name string `gorm:"varchar(60); not null; unique"`
	//Artists []Artist
}
