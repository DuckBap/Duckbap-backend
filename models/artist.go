package models

type Artist struct {
	EnterID uint `gorm:"not null"`
	Name string `gorm:"varchar(60); not null"`
	ImgUrl string `gorm:"varchar(255); unique; not null"`
}
