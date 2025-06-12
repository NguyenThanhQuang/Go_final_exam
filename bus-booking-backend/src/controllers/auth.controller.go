package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"

	"github.com/Go_final_exam/bus-booking-backend/src/services"
)

var validate = validator.New()


// @Summary Đăng ký tài khoản người dùng mới
// @Description Tạo một tài khoản mới cho người dùng với email, số điện thoại, tên và mật khẩu.
// @Tags Authentication
// @Accept  json
// @Produce  json
// @Param   user body services.RegisterInput true "Thông tin đăng ký của người dùng"
// @Success 201 {object} map[string]interface{} "Trả về thông tin người dùng đã tạo"
// @Failure 400 {object} map[string]string "Dữ liệu đầu vào không hợp lệ"
// @Failure 409 {object} map[string]string "Email đã được sử dụng"
// @Router /auth/register [post]
func RegisterController(c *gin.Context) {
	var input services.RegisterInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"lỗi": "Dữ liệu đầu vào không hợp lệ"})
		return
	}

	if err := validate.Struct(input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"lỗi": err.Error()})
		return
	}

	user, err := services.Register(input)
	if err != nil {
		c.JSON(http.StatusConflict, gin.H{"lỗi": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"thông báo": "Đăng ký tài khoản thành công!",
		"dữ liệu": gin.H{
			"id":        user.ID,
			"name":      user.Name,
			"email":     user.Email,
			"phone":     user.Phone,
			"createdAt": user.CreatedAt,
		},
	})
}

// @Summary Đăng nhập vào hệ thống
// @Description Xác thực người dùng bằng email và mật khẩu, trả về một token JWT.
// @Tags Authentication
// @Accept  json
// @Produce  json
// @Param   credentials body services.LoginInput true "Thông tin đăng nhập"
// @Success 200 {object} map[string]string "Đăng nhập thành công, trả về token"
// @Failure 400 {object} map[string]string "Dữ liệu đầu vào không hợp lệ"
// @Failure 401 {object} map[string]string "Email hoặc mật khẩu không chính xác"
// @Router /auth/login [post]
func LoginController(c *gin.Context) {
	var input services.LoginInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"lỗi": "Dữ liệu đầu vào không hợp lệ"})
		return
	}

	if err := validate.Struct(input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"lỗi": err.Error()})
		return
	}

	token, err := services.Login(input)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"lỗi": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"thông báo": "Đăng nhập thành công!",
		"token":     token,
	})
}