package models

import (
	"gorm.io/gorm"
)


type User struct {
	gorm.Model
	UserName string `gorm:"type:varchar(50); not null; unique" json:"userName"`
	Password string `gorm:"type:varchar(255); not null" json:"password"`
	Email string `gorm:"type:varchar(100); not null; unique" json:"email"`
	NickName string `gorm:"type:varchar(50); not null; unique" json:"nickName"`
	FavoriteArtist uint `gorm:"not null" json:"favoriteArtist"`
	Fundings []Funding `gorm:"foreignKey:SellerID"`
	Receipts []Receipt `gorm:"foreignKey:SellerID"`
	Artist Artist `gorm:"foreignKey:FavoriteArtist"`
}
