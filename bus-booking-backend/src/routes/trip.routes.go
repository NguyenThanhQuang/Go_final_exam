package routes

import (
	"github.com/Go_final_exam/bus-booking-backend/src/controllers"
	"github.com/gin-gonic/gin"
)

func TripRoutes(router *gin.RouterGroup) {
	tripGroup := router.Group("/trips")
	{
		tripGroup.GET("", controllers.SearchTripsController)
		
		tripGroup.GET("/:tripId", controllers.GetTripDetailsController)
	}
}