package middleware

import (
	"bus_booking_api/internal/config"
	"bus_booking_api/internal/utils"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Yêu cầu chưa được xác thực (thiếu token)."})
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Định dạng token không hợp lệ (phải là 'Bearer <token>')."})
			return
		}
		tokenString := parts[1]

		claims, err := utils.ValidateToken(tokenString, cfg)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Token không hợp lệ hoặc đã hết hạn: " + err.Error()})
			return
		}

		c.Set("userID", claims.UserID)
		c.Set("userEmail", claims.Email)
		c.Set("userRole", claims.Role)
		if claims.CompanyID != "" {
			c.Set("userCompanyID", claims.CompanyID)
		}

		c.Next() 
	}
}

func RoleAuthMiddleware(allowedRoles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userRole, exists := c.Get("userRole")
		if !exists {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Không thể xác định vai trò người dùng."})
			return
		}

		roleStr, ok := userRole.(string)
		if !ok {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Vai trò người dùng không hợp lệ."})
			return
		}

		isAllowed := false
		for _, allowedRole := range allowedRoles {
			if roleStr == allowedRole {
				isAllowed = true
				break
			}
		}

		if !isAllowed {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Bạn không có quyền truy cập tài nguyên này."})
			return
		}
		c.Next()
	}
}