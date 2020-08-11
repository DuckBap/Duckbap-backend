package models

import (
	"gorm.io/gorm"
)


type User struct {
	gorm.Model
	UserName string `gorm:"varchar(50); not null; unique" json:"userName"`
	Password string `gorm:"varchar(255); not null" json:"password"`
	Email string `gorm:"varchar(100); not null; unique" json:"email"`
	NickName string `gorm:"varchar(50); not null; unique" json:"nickName"`
	FavoriteArtist uint `gorm:"not null" json:"favoriteArtist"`
	Fundings []Funding `gorm:"foreignKey:SellerID"`
	Receipts []Receipt `gorm:"foreignKey:SellerID"`
	Artist Artist `gorm:"foreignKey:FavoriteArtist"`
}
