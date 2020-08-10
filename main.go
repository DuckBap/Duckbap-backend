package main

import (
	"github.com/DuckBap/Duckbap-backend/configs"
	"github.com/DuckBap/Duckbap-backend/models"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
)

func main() {
	var err error
	r := gin.New()

	//newLogger := logger.New(
	//	log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
	//	logger.Config{
	//		SlowThreshold: time.Second,   // Slow SQL threshold
	//		LogLevel:      logger.Info, // Log level
	//		Colorful:      false,         // Disable color
	//	},
	//)

	configs.DB, err = gorm.Open(mysql.Open(configs.DbURL(configs.BuildDBConfig())), &gorm.Config{
		//Logger: newLogger,
	})
	if err != nil {
		log.Println(err)
	}
	configs.DB.AutoMigrate(&models.User{}, &models.Funding{}, &models.FundingImg{},
		&models.Artist{}, &models.Receipt{}, &models.Entertainment{})

	//rGroup := r.Group("/")
	//routers.SetFundingUrls(rGroup.Group("/fundings"))

	r.Run(":8080")
}