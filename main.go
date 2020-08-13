package main

import (
	"fmt"
	"github.com/DuckBap/Duckbap-backend/configs"
	"github.com/DuckBap/Duckbap-backend/docs"
	"github.com/DuckBap/Duckbap-backend/models"
	"github.com/DuckBap/Duckbap-backend/routers"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
	"os"
)

func main() {
	var err error
	r := gin.Default()

	//개발용 로컬 디비 오픈
	//configs.DB, err = gorm.Open(mysql.Open(configs.DbURL(configs.BuildDBConfig())), &gorm.Config{
	//	//Logger: newLogger,
	//})

	//배포용 리모트 디비 오픈
	configs.DB, err = initSocketConnectionPool()

	if err != nil {
		log.Println(err)
	}

	//개발용 마이그레이트
	//configs.DB.AutoMigrate(&models.User{}, &models.Funding{}, &models.FundingImg{},
	//	&models.Artist{}, &models.Receipt{}, &models.Entertainment{})

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

	r.Run()
}

//배포용
func mustGetenv(k string) string {
	v := os.Getenv(k)
	if v == "" {
		log.Fatalf("Warning: %s environment variable not set.\n", k)
	}
	return v
}
//배포용
func initSocketConnectionPool() (*gorm.DB, error) {
	var (
		dbUser                 = mustGetenv("DB_USER")
		dbPwd                  = mustGetenv("DB_PASS")
		instanceConnectionName = mustGetenv("INSTANCE_CONNECTION_NAME")
		dbName                 = mustGetenv("DB_NAME")
	)
	socketDir := "/cloudsql"
	var dbURI string
	dbURI = fmt.Sprintf("%s:%s@unix(/%s/%s)/%s?parseTime=true", dbUser, dbPwd, socketDir, instanceConnectionName, dbName)
	dbPool, err := gorm.Open(mysql.Open(dbURI), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("sql.Open: %v", err)
	}
	dbPool.AutoMigrate(&models.User{}, &models.Funding{}, &models.FundingImg{},
		&models.Artist{}, &models.Receipt{}, &models.Entertainment{})
	//configureConnectionPool(dbPool)
	return dbPool, nil
}
