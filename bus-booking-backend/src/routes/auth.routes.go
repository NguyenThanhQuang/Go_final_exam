package routes

import (
	"github.com/Go_final_exam/bus-booking-backend/src/controllers"
	"github.com/gin-gonic/gin"
)

func AuthRoutes(router *gin.RouterGroup) {
	authGroup := router.Group("/auth")
	{
		authGroup.POST("/register", controllers.RegisterController)
		authGroup.POST("/login", controllers.LoginController)
	}
}