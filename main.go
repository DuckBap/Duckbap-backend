package main

import (
	"github.com/DuckBap/Duckbap-backend/configs"
	"github.com/DuckBap/Duckbap-backend/models"
	"github.com/DuckBap/Duckbap-backend/routers"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm/logger"
	"os"
	"time"

	//"github.com/DuckBap/duckBap/routers"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
)

func main() {
	var err error
	r := gin.New()

	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold: time.Second,   // Slow SQL threshold
			LogLevel:      logger.Info, // Log level
			Colorful:      false,         // Disable color
		},
	)
	configs.DB, err = gorm.Open(mysql.Open(configs.DbURL(configs.BuildDBConfig())), &gorm.Config{
		Logger: newLogger,
	})
	if err != nil {
		log.Println(err)
	}
	configs.DB.AutoMigrate(&models.User{}, &models.Funding{}, &models.FundingImg{},&models.Artist{}, &models.Receipt{},  &models.Entertainment{})
	//configs.DB.AutoMigrate(&models.Entertainment{})
	//configs.DB.AutoMigrate(&models.Artist{})//, &models.Entertainment{})
	//configs.DB.AutoMigrate(&models.User{})
	//configs.DB.AutoMigrate(&models.Funding{})
	//configs.DB.AutoMigrate(&models.Receipt{}, &models.FundingImg{})

	rGroup := r.Group("/")
	routers.SetFundingRouter(rGroup.Group("/fundings"))
	r.Run(":8080")
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
