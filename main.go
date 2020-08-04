package main

import (
	"github.com/DuckBap/duckBap/configs"
	"github.com/DuckBap/duckBap/models"
	//"github.com/DuckBap/duckBap/routers"
	//"github.com/gin-gonic/gin"
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
	configs.DB.AutoMigrate(&models.FundingImg{},&models.Artist{})
	//r := routers.SetupRouter()
	//r.Run()
}
