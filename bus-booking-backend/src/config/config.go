package config

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Config struct {
	Port                string
	MongoURI            string
	MongoDatabaseName   string
	JwtSecretKey        string
	JwtExpirationInHours string
}

var DB *mongo.Database

func LoadConfig() (Config, error) {
	err := godotenv.Load()
	if err != nil {
		log.Println("Cảnh báo: Không tìm thấy file .env, đang sử dụng biến môi trường hệ thống.")
	}

	cfg := Config{
		Port:                os.Getenv("PORT"),
		MongoURI:            os.Getenv("MONGODB_URI"),
		MongoDatabaseName:   os.Getenv("MONGODB_DATABASE_NAME"),
		JwtSecretKey:        os.Getenv("JWT_SECRET_KEY"),
		JwtExpirationInHours: os.Getenv("JWT_EXPIRATION_HOURS"),
	}

	if cfg.Port == "" || cfg.MongoURI == "" || cfg.MongoDatabaseName == "" || cfg.JwtSecretKey == "" {
		return Config{}, fmt.Errorf("lỗi: một hoặc nhiều biến môi trường quan trọng chưa được thiết lập")
	}

	return cfg, nil
}

func ConnectDB(cfg Config) {
	clientOptions := options.Client().ApplyURI(cfg.MongoURI)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatalf("Không thể kết nối đến MongoDB: %v", err)
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatalf("Không thể ping MongoDB: %v", err)
	}

	log.Println("Đã kết nối thành công đến MongoDB!")
	DB = client.Database(cfg.MongoDatabaseName)
}