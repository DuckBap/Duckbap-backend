package models

import (
	"gorm.io/gorm"
)


type User struct {
	gorm.Model
	UserName string `gorm:"varchar(50); not null; unique"`
	Password string `gorm:"varchar(255); not null"`
	Email string `gorm:"varchar(100); not null; unique"`
	NickName string `gorm:"varchar(50); not null; unique"`
	FavoriteArtist uint `gorm:"not null"`
	Fundings []Funding `gorm:"foreignKey:SellerID"`
	Receipts []Receipt `gorm:"foreignKey:SellerID"`
	Artist Artist `gorm:"foreignKey:FavoriteArtist"`
}
