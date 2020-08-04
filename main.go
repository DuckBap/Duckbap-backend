package main

import (
	"github.com/DuckBap/duckBap/configs"
	"github.com/DuckBap/duckBap/models"
	"github.com/gin-gonic/gin"
	//"github.com/DuckBap/duckBap/routers"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
)

func main() {
	var err error

	configs.DB, err = gorm.Open(mysql.Open(configs.DbURL(configs.BuildDBConfig())), &gorm.Config{})
	if err != nil {
		log.Println(err)
	}
	//configs.DB.AutoMigrate(&models.User{}, &models.Funding{}, &models.FundingImg{},&models.Artist{}, &models.Receipt{},  &models.Entertainment{})
	configs.DB.AutoMigrate(&models.Entertainment{})

	r := gin.Default()
	//r.GET("/", Test)
	r.Run()
}
//
//func Test(c *gin.Context){
//
//	var ff models.Funding
//	ff.StartDate = time.Now()
//	configs.DB.Create(&ff)
//
//	configs.DB.Find(&ff, "id = ?", 1)
//}
