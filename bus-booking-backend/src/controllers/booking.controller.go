package controllers

import (
	"net/http"

	"github.com/Go_final_exam/bus-booking-backend/src/services"
	"github.com/gin-gonic/gin"
)

// @Summary Tạo một booking mới (Giữ chỗ)
// @Description Giữ chỗ cho người dùng đã đăng nhập. Yêu cầu token xác thực.
// @Tags Bookings
// @Accept  json
// @Produce  json
// @Security BearerAuth
// @Param   booking body services.CreateBookingInput true "Thông tin để giữ chỗ"
// @Success 201 {object} map[string]interface{} "Giữ chỗ thành công, trả về thông tin booking"
// @Failure 400 {object} map[string]string "Dữ liệu đầu vào không hợp lệ"
// @Failure 401 {object} map[string]string "Yêu cầu token xác thực"
// @Failure 404 {object} map[string]string "Không tìm thấy chuyến đi"
// @Failure 409 {object} map[string]string "Ghế đã được người khác chọn (Conflict)"
// @Router /bookings [post]
func CreateBookingController(c *gin.Context) {
	var input services.CreateBookingInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"lỗi": "Dữ liệu đầu vào không hợp lệ"})
		return
	}

	if err := validate.Struct(input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"lỗi": err.Error()})
		return
	}

	userID, exists := c.Get("userId")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"lỗi": "Không tìm thấy thông tin người dùng. Vui lòng đăng nhập lại."})
		return
	}
	userIDStr, ok := userID.(string)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"lỗi": "Lỗi định dạng ID người dùng."})
		return
	}

	booking, err := services.CreateBooking(input, userIDStr)
	if err != nil {
		if err.Error() == "không tìm thấy chuyến đi" {
			c.JSON(http.StatusNotFound, gin.H{"lỗi": err.Error()})
			return
		}
		if err.Error() == "một hoặc nhiều ghế đã được người khác chọn. Vui lòng thử lại" {
			c.JSON(http.StatusConflict, gin.H{"lỗi": err.Error()})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"lỗi": err.Error()}) 
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"thông báo": "Giữ chỗ thành công! Vui lòng tiến hành thanh toán.",
		"dữ liệu":   booking,
	})
}