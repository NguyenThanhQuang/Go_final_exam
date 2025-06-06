package handlers

import (
	"bus_booking_api/internal/models"
	"bus_booking_api/internal/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	authService services.AuthService
}

func NewAuthHandler(authService services.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

// Register godoc
// @Summary Đăng ký tài khoản người dùng mới
// @Description Tạo một tài khoản người dùng mới với email, số điện thoại, mật khẩu và tên.
// @Tags auth
// @Accept  json
// @Produce  json
// @Param   user body models.RegisterRequest true "Thông tin đăng ký"
// @Success 201 {object} models.User "Tạo tài khoản thành công"
// @Failure 400 {object} map[string]string "Dữ liệu không hợp lệ"
// @Failure 500 {object} map[string]string "Lỗi máy chủ nội bộ"
// @Router /auth/register [post]
func (h *AuthHandler) Register(c *gin.Context) {
	var req models.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Dữ liệu không hợp lệ: " + err.Error()})
		return
	}

	user, err := h.authService.Register(c.Request.Context(), req.Name, req.Email, req.Phone, req.Password)
	if err != nil {
		// Phân biệt lỗi do người dùng hay lỗi server
		if err.Error() == "email đã được sử dụng" || err.Error() == "số điện thoại đã được sử dụng" ||
			err.Error() == "tên, email, số điện thoại và mật khẩu không được để trống" ||
			err.Error() == "mật khẩu phải có ít nhất 6 ký tự" ||
			err.Error() == "địa chỉ email không hợp lệ" {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Không thể đăng ký tài khoản: " + err.Error()})
		}
		return
	}

	c.JSON(http.StatusCreated, user)
}

// Login godoc
// @Summary Đăng nhập vào hệ thống
// @Description Đăng nhập bằng email và mật khẩu, trả về thông tin người dùng và JWT token.
// @Tags auth
// @Accept  json
// @Produce  json
// @Param   credentials body models.LoginRequest true "Thông tin đăng nhập"
// @Success 200 {object} models.LoginResponse "Đăng nhập thành công"
// @Failure 400 {object} map[string]string "Dữ liệu không hợp lệ"
// @Failure 401 {object} map[string]string "Thông tin đăng nhập không chính xác"
// @Failure 500 {object} map[string]string "Lỗi máy chủ nội bộ"
// @Router /auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	var req models.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Dữ liệu không hợp lệ: " + err.Error()})
		return
	}

	user, token, err := h.authService.Login(c.Request.Context(), req.Email, req.Password)
	if err != nil {
		if err.Error() == "email hoặc mật khẩu không chính xác" ||
			err.Error() == "email và mật khẩu không được để trống" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Lỗi đăng nhập: " + err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, models.LoginResponse{Token: token, User: *user})
}