package middlewares

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"

	"github.com/Go_final_exam/bus-booking-backend/src/config"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		cfg, err := config.LoadConfig()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"lỗi": "Không thể tải cấu hình hệ thống"})
			c.Abort()
			return
		}
		jwtSecretKey := cfg.JwtSecretKey

		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"lỗi": "Yêu cầu cần có token xác thực"})
			c.Abort()
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"lỗi": "Định dạng token không hợp lệ"})
			c.Abort()
			return
		}

		tokenString := parts[1]

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, http.ErrAbortHandler
			}
			return []byte(jwtSecretKey), nil
		})

		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"lỗi": "Token không hợp lệ hoặc đã hết hạn"})
			c.Abort()
			return
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			c.Set("userId", claims["userId"])
		} else {
			c.JSON(http.StatusUnauthorized, gin.H{"lỗi": "Token không hợp lệ"})
			c.Abort()
			return
		}

		c.Next()
	}
}