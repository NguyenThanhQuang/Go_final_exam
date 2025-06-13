package routes

import (
	"github.com/Go_final_exam/bus-booking-backend/src/controllers"
	"github.com/Go_final_exam/bus-booking-backend/src/middlewares"
	"github.com/gin-gonic/gin"
)

func BookingRoutes(router *gin.RouterGroup) {
	bookingGroup := router.Group("/bookings")
  	bookingGroup.Use(middlewares.AuthMiddleware())
  	{
  		bookingGroup.POST("", controllers.CreateBookingController)
  		bookingGroup.GET("/:bookingId", controllers.GetBookingDetailsController) 
  		bookingGroup.GET("/my", controllers.GetMyBookingsController)         
  	}
}