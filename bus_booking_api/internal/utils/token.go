package utils

import (
	"bus_booking_api/internal/config"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)
type JWTCustomClaims struct {
	UserID    string `json:"userId"`
	Email     string `json:"email"`
	Role      string `json:"role"`
	CompanyID string `json:"companyId,omitempty"` 
	jwt.RegisteredClaims
}

func GenerateToken(userID, email, role, companyID string, cfg *config.Config) (string, error) {
	claims := &JWTCustomClaims{
		UserID:    userID,
		Email:     email,
		Role:      role,
		CompanyID: companyID, 
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(cfg.JWTExpiration)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "bus-booking-api", 
			Subject:   userID,           
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(cfg.JWTSecretKey))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func ValidateToken(tokenString string, cfg *config.Config) (*JWTCustomClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &JWTCustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New(" thuật toán ký không mong muốn")
		}
		return []byte(cfg.JWTSecretKey), nil
	})

	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, errors.New("token đã hết hạn")
		}
		if errors.Is(err, jwt.ErrTokenNotValidYet) {
			return nil, errors.New("token chưa hợp lệ")
		}
		return nil, errors.New("token không hợp lệ")
	}

	if claims, ok := token.Claims.(*JWTCustomClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("token không hợp lệ")
}