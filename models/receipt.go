package models

import "gorm.io/gorm"

type Receipt struct {
	gorm.Model
	SellerID uint `gorm:"not null" json:"sellerId"`
	BuyerID uint `gorm:"not null" json:"buyerId"`
	FundingID uint `gorm:"not null" json:"fundingId"`
	Amount uint `gorm:"not null" json:"amount"`
	Funding Funding
	User User `gorm:"foreignKey:BuyerID"`
}
