package config

import (
	"log"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	Port                string
	MongoDBURI          string
	MongoDBDatabaseName string
	JWTSecretKey        string
	JWTExpiration       time.Duration
}

var AppConfig *Config

func LoadConfig() (*Config, error) {
	err := godotenv.Load()
	if err != nil {
		log.Println("Thông báo: Không tìm thấy file .env, sẽ sử dụng biến môi trường hệ thống (nếu có).")
	}

	port := getEnv("PORT", "8080") 
	mongoURI := getEnv("MONGODB_URI", "mongodb://localhost:27017")
	mongoDBName := getEnv("MONGODB_DATABASE_NAME", "bus_booking_db")
	jwtSecret := getEnv("JWT_SECRET_KEY", "your-secret-key-default") 
	jwtExpHoursStr := getEnv("JWT_EXPIRATION_HOURS", "1")

	jwtExpHours, err := strconv.Atoi(jwtExpHoursStr)
	if err != nil {
		log.Printf("Cảnh báo: JWT_EXPIRATION_HOURS không hợp lệ ('%s'), sử dụng giá trị mặc định là 1 giờ.\n", jwtExpHoursStr)
		jwtExpHours = 1
	}

	if jwtSecret == "your-secret-key-default" {
		log.Println("CẢNH BÁO: JWT_SECRET_KEY đang sử dụng giá trị mặc định. Vui lòng đặt giá trị an toàn trong file .env hoặc biến môi trường.")
	}

	AppConfig = &Config{
		Port:                port,
		MongoDBURI:          mongoURI,
		MongoDBDatabaseName: mongoDBName,
		JWTSecretKey:        jwtSecret,
		JWTExpiration:       time.Hour * time.Duration(jwtExpHours),
	}
	log.Println("Đã tải cấu hình ứng dụng thành công.")
	return AppConfig, nil
}

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}

func GetDBCollectionName(name string) string {
	return name
}

const (
	UsersCollection    = "users"
	CompaniesCollection = "companies"
	VehiclesCollection = "vehicles"
	TripsCollection    = "trips"
	BookingsCollection = "bookings"
)