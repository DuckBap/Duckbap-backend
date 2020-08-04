package models

import "gorm.io/gorm"

type Receipt struct {
	gorm.Model
	SellerID uint `gorm:"not null"`
	BuyerID uint `gorm:"not null"`
	Amount uint `gorm:"not null"`
	FundingID uint `gorm:"not null"`
	Funding Funding
	User User `gorm:"foreignKey:BuyerID"`
}
