package main

import (
	"github.com/DuckBap/Duckbap-backend/configs"
	"github.com/DuckBap/Duckbap-backend/models"
	"github.com/DuckBap/Duckbap-backend/routers"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
)

func main() {
	var err error
	r := gin.New()
	configs.DB, err = gorm.Open(mysql.Open(configs.DbURL(configs.BuildDBConfig())), &gorm.Config{})
	if err != nil {
		log.Println(err)
	}
	configs.DB.AutoMigrate(&models.User{}, &models.Funding{}, &models.FundingImg{},&models.Artist{}, &models.Receipt{},  &models.Entertainment{})
	routerGroup := r.Group("/")
	routers.SignUpRouter(routerGroup)
	routers.SetUserRouters(routerGroup)
	r.Run(":8080")
}

