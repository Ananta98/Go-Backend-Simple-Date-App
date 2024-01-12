package routes

import (
	"backend-dating-app/controllers"
	"backend-dating-app/middlewares"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func SetupRouter(db *gorm.DB) *gin.Engine {
	r := gin.Default()
	r.Use(func(ctx *gin.Context) {
		ctx.Set("db", db)
	})

	authRoute := r.Group("auth")
	authRoute.POST("/register", controllers.Register)
	authRoute.POST("/login", controllers.Login)

	dateMiddlewareRoute := r.Group("date")
	dateMiddlewareRoute.Use(middlewares.JWTAuthMiddleware())
	dateMiddlewareRoute.GET("/list-matching", controllers.ListMatching)

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	return r
}
