package main

import (
	"log"
	"os"

	configuration "github.com/shravanasati/scopex-go-assignment/configuration"
	cronjob "github.com/shravanasati/scopex-go-assignment/cronjob"
	docs "github.com/shravanasati/scopex-go-assignment/docs"
	router "github.com/shravanasati/scopex-go-assignment/router"
	util "github.com/shravanasati/scopex-go-assignment/util"

	"github.com/spf13/viper"
	swgFiles "github.com/swaggo/files"
	swgGin "github.com/swaggo/gin-swagger"
)

func init() {

	os.Setenv("APP_ENVIRONMENT", "STAGING")

	// read config environment
	configuration.ReadConfig()

	util.Pool = util.SetupRedisJWT()

}

// @securityDefinitions.apikey bearerAuth
// @in header
// @name Authorization
func main() {

	var err error

	// Setup database
	configuration.DB, err = configuration.SetupDB()
	if err != nil {
		log.Fatal(err)
	}
	defer configuration.DB.Close()

	// Start Cron Jobs
	cronjob.InitCron()

	port := viper.GetString("PORT")

	docs.SwaggerInfo.Title = "Swagger Service API"
	docs.SwaggerInfo.Description = "This is service API documentation."
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.Host = "localhost:" + port
	docs.SwaggerInfo.BasePath = "/api"
	docs.SwaggerInfo.Schemes = []string{"http", "https"}

	// Setup router
	router := router.NewRoutes()
	url := swgGin.URL("http://localhost:" + port + "/swagger/doc.json")
	router.GET("/swagger/*any", swgGin.WrapHandler(swgFiles.Handler, url))

	log.Fatal(router.Run(":" + port))
}
