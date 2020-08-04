package models

import (
	"gorm.io/gorm"
	"time"
)

type Funding struct {
	gorm.Model
	SellerID uint `gorm:"not null"`
	Name string `gorm:"varchar(150);not null;"`
	Price uint `gorm:"not null"`
	TargetAmount uint `gorm:"not null"`
	StartDate time.Time
	EndDate time.Time
}

type FundingImg struct {
	gorm.Model
	FundingID uint `gorm:"not null"`
	Url string `gorm:"type:varchar(255); unique; not null"`
	IsTitle bool `gorm:"default:false; not null"`
	Order uint8 `gorm:"not null"`
}
