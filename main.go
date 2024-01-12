package main

import (
	"backend-dating-app/config"
	"backend-dating-app/docs"
	"backend-dating-app/routes"
	"backend-dating-app/utils"

	"github.com/joho/godotenv"
)

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @termsOfService http://swagger.io/terms/
func main() {
	err := godotenv.Load()
	if err != nil {
		panic(err.Error())
	}
	// programmatically set swagger info
	docs.SwaggerInfo.Title = "Swagger Date App API"
	docs.SwaggerInfo.Description = "API Documentation for Date App."
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.Host = utils.GetEnv("SWAGGER_HOST", "localhost:8080")
	docs.SwaggerInfo.Schemes = []string{"http", "https"}

	db := config.ConnectDatabase()
	sqlDB, err := db.DB()
	if err != nil {
		panic(err.Error())
	}
	defer sqlDB.Close()

	r := routes.SetupRouter(db)
	r.Run()
}
