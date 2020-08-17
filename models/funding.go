package models

import (
	"gorm.io/gorm"
	"time"
)

type Funding struct {
	gorm.Model
	SellerID uint `gorm:"not null" json:"sellerId"`
	Name string `gorm:"type:varchar(150);not null;" json:"name"`
	Price uint `gorm:"not null" json:"price"`
	TargetAmount uint `gorm:"not null" json:"targetAmount"`
	StartDate time.Time `gorm:"type:date; not null" json:"startDate"`
	EndDate time.Time `gorm:"type:date; not null" json:"endDate"`
	MainImgUrl string `gorm:"type:varchar(512); unique; not null" json:"mainImgUrl"`
	ArtistID uint `gorm:"not null" json:"artistId"`
	SalesAmount uint `gorm:"default:0" json:"salesAmount"`
	FundingImgs []FundingImg
	Receipts []Receipt
	Artist Artist
}

type FundingImg struct {
	gorm.Model
	FundingID uint `gorm:"not null" json:"fundingId"`
	Url string `gorm:"type:varchar(512); unique; not null" json:"url"`
	IsTitle bool `gorm:"default:false; not null" json:"isTitle"`
	Order uint8 `gorm:"not null" json:"order"`
}
