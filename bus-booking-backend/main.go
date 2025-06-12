package main

import (
	"log"

	"github.com/Go_final_exam/bus-booking-backend/docs"
	"github.com/Go_final_exam/bus-booking-backend/src/config"
	"github.com/Go_final_exam/bus-booking-backend/src/routes"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title API Dịch vụ Đặt vé xe
// @version 1.0
// @description Đây là tài liệu API cho ứng dụng Backend đặt vé xe viết bằng Go.
// @description Final Assignment - VTC Academy - Build Backend with Golang.

// @contact.name API Support
// @contact.email your.email@example.com

// @host localhost:8080
// @BasePath /api/v1

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Nhập token JWT với tiền tố 'Bearer '. Ví dụ: "Bearer {token}"
func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Không thể tải cấu hình: %v", err)
	}

	config.ConnectDB(cfg)

	docs.SwaggerInfo.Title = "API Dịch vụ Đặt vé xe"
	docs.SwaggerInfo.Description = "Đây là tài liệu API cho ứng dụng Backend đặt vé xe viết bằng Go."
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.Host = cfg.Port 
	docs.SwaggerInfo.BasePath = "/api/v1"
	docs.SwaggerInfo.Schemes = []string{"http", "httpshttps"}

	gin.SetMode(gin.ReleaseMode)

	router := gin.New()

	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	router.SetTrustedProxies([]string{"127.0.0.1"})

	router.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	api := router.Group("/api/v1")
	{
		routes.AuthRoutes(api)
		routes.TripRoutes(api)
		routes.BookingRoutes(api)
	}
	
	log.Printf("Server đang chạy trên cổng %s", cfg.Port)
	if err := router.Run(":" + cfg.Port); err != nil {
		log.Fatalf("Lỗi khi khởi động server: %v", err)
	}
}