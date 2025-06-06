package main

import (
	"bus_booking_api/internal/config"
	"bus_booking_api/internal/handlers"
	"bus_booking_api/internal/middleware"
	"bus_booking_api/internal/repositories"
	"bus_booking_api/internal/seed"
	"bus_booking_api/internal/services"
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"

	_ "bus_booking_api/docs"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

var mongoClient *mongo.Client

// InitMongoDB initializes the MongoDB connection
func InitMongoDB(cfg *config.Config) (*mongo.Database, error) {
	log.Println("Đang kết nối tới MongoDB...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	clientOptions := options.Client().ApplyURI(cfg.MongoDBURI)
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Printf("Lỗi khi kết nối tới MongoDB: %v\n", err)
		return nil, err
	}

	// Ping the primary
	if err := client.Ping(ctx, readpref.Primary()); err != nil {
		log.Printf("Lỗi khi ping MongoDB: %v\n", err)
		return nil, err
	}

	mongoClient = client // Store the client globally or pass it around
	db := client.Database(cfg.MongoDBDatabaseName)
	log.Printf("Đã kết nối thành công tới MongoDB, cơ sở dữ liệu: %s\n", cfg.MongoDBDatabaseName)
	return db, nil
}

// CloseMongoDBConnection closes the MongoDB connection
func CloseMongoDBConnection() {
	if mongoClient != nil {
		log.Println("Đang đóng kết nối MongoDB...")
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := mongoClient.Disconnect(ctx); err != nil {
			log.Printf("Lỗi khi đóng kết nối MongoDB: %v\n", err)
		} else {
			log.Println("Đã đóng kết nối MongoDB.")
		}
	}
}

// @title Bus Booking API
// @version 1.0
// @description API cho hệ thống đặt vé xe khách.
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8080
// @BasePath /api/v1

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
func main() {
	// 1. Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Lỗi khi tải cấu hình: %v", err)
	}

	// 2. Initialize MongoDB connection
	db, err := InitMongoDB(cfg)
	if err != nil {
		log.Fatalf("Không thể khởi tạo kết nối MongoDB: %v", err)
	}
	defer CloseMongoDBConnection()

	// 3. Initialize Repositories
	userRepo := repositories.NewMongoUserRepository(db, config.UsersCollection)
	companyRepo := repositories.NewMongoCompanyRepository(db, config.CompaniesCollection)
	vehicleRepo := repositories.NewMongoVehicleRepository(db, config.VehiclesCollection)
	tripRepo := repositories.NewMongoTripRepository(db, config.TripsCollection)
	bookingRepo := repositories.NewMongoBookingRepository(db, config.BookingsCollection)
	log.Println("Tất cả repositories đã được khởi tạo.")

	// 4. Initialize Services
	authSvc := services.NewAuthService(userRepo, cfg)
	tripSvc := services.NewTripService(tripRepo, companyRepo, vehicleRepo) // Truyền companyRepo và vehicleRepo
	bookingSvc := services.NewBookingService(bookingRepo, tripRepo, userRepo)
	log.Println("Tất cả services đã được khởi tạo.")

	// --- SEED DATA ---
	// Gọi hàm RunSeeders từ package seed
	seed.RunSeeders(cfg, userRepo, companyRepo, vehicleRepo, tripRepo, bookingRepo, tripSvc)
	// --- END SEED DATA ---

	// 5. Initialize Handlers
	authHandler := handlers.NewAuthHandler(authSvc)
	tripHandler := handlers.NewTripHandler(tripSvc)
	bookingHandler := handlers.NewBookingHandler(bookingSvc)
	log.Println("Tất cả handlers đã được khởi tạo.")

	// 6. Initialize Gin Router
	router := gin.Default()

	// CORS Middleware
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// API v1 routes
	apiV1 := router.Group("/api/v1")
	{
		// Auth routes
		authRoutes := apiV1.Group("/auth")
		{
			authRoutes.POST("/register", authHandler.Register)
			authRoutes.POST("/login", authHandler.Login)
		}

		// Trip routes
		tripRoutes := apiV1.Group("/trips")
		{
			tripRoutes.GET("/search", tripHandler.SearchTrips)
			tripRoutes.GET("/:tripId", tripHandler.GetTripDetails)
		}

		// Booking routes (yêu cầu xác thực)
		bookingRoutes := apiV1.Group("/bookings")
		bookingRoutes.Use(middleware.AuthMiddleware(cfg)) // Áp dụng middleware xác thực
		{
			bookingRoutes.POST("", bookingHandler.CreateBooking)
			bookingRoutes.POST("/payment/mock", bookingHandler.MockPayment)
			bookingRoutes.GET("/:bookingId", bookingHandler.GetBookingDetails)
		}
	}

	// Swagger UI
	// Đường dẫn để truy cập Swagger UI: http://localhost:PORT/swagger/index.html
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// 7. Start HTTP server
	server := &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: router,
	}

	go func() {
		log.Printf("Đang lắng nghe và phục vụ HTTP trên cổng %s\n", cfg.Port)
		log.Printf("Swagger UI có sẵn tại: http://localhost:%s/swagger/index.html", cfg.Port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Lỗi lắng nghe: %s\n", err)
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Đang tắt server...")

	ctxTimeout, cancelTimeout := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelTimeout()
	if err := server.Shutdown(ctxTimeout); err != nil {
		log.Fatalf("Lỗi khi tắt server: %v", err)
	}

	log.Println("Server đã tắt.")
}
