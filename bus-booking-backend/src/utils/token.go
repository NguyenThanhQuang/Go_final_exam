package utils

import (
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/Go_final_exam/bus-booking-backend/src/config"
)

func GenerateToken(userID primitive.ObjectID) (string, error) {
	cfg, err := config.LoadConfig()
	if err != nil {
		return "", err
	}

	expirationHours, err := strconv.Atoi(cfg.JwtExpirationInHours)
	if err != nil {
		return "", err
	}

	claims := jwt.MapClaims{
		"userId": userID.Hex(),                                     
		"exp":    time.Now().Add(time.Hour * time.Duration(expirationHours)).Unix(),
		"iat":    time.Now().Unix(),                                               
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString([]byte(cfg.JwtSecretKey))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}