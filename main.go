package main

import (
	"github.com/DuckBap/Duckbap-backend/configs"
	"github.com/DuckBap/Duckbap-backend/docs"
	"github.com/DuckBap/Duckbap-backend/models"
	"github.com/DuckBap/Duckbap-backend/routers"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
)

func main() {
	var err error
	r := gin.Default()
	r.Use(cors.Default())

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

	rGroup := r.Group("/v1")
	routers.SignUpRouter(rGroup.Group("/accounts"))
	routers.SetUserRouters(rGroup.Group("/users"))
	routers.SetArtistRouter(rGroup.Group("/artists"))
	routers.SetFundingRouter(rGroup.Group("/fundings"))
	routers.SetEnterRouter(rGroup.Group("/ents"))
	routers.SetSwagger(rGroup.Group("/swagger"))

	docs.SwaggerInfo.Title = "duckBap API"
	docs.SwaggerInfo.Description = ""
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.Host = "localhost:8080"
	docs.SwaggerInfo.BasePath = "/v1"

	r.Run(":8080")
}
