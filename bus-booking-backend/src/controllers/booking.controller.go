package controllers

import (
	"net/http"
	"strings"

	"github.com/Go_final_exam/bus-booking-backend/src/models"
	"github.com/Go_final_exam/bus-booking-backend/src/services"
	"github.com/gin-gonic/gin"
)

// @Summary Tạo một booking mới (Giữ chỗ)
// @Description Giữ chỗ cho người dùng đã đăng nhập. Yêu cầu token xác thực.
// @Tags Bookings
// @Accept  json
// @Produce  json
// @Security BearerAuth
// @Param   booking body services.CreateBookingInput true "Thông tin để giữ chỗ (tripId, seatNumbers)"
// @Success 201 {object} map[string]interface{} "Giữ chỗ thành công. Body: {thông báo: string, dữ_liệu: models.Booking}"
// @Failure 400 {object} map[string]string "Dữ liệu đầu vào không hợp lệ hoặc lỗi xử lý khác"
// @Failure 401 {object} map[string]string "Yêu cầu token xác thực hoặc không tìm thấy thông tin người dùng"
// @Failure 404 {object} map[string]string "Không tìm thấy chuyến đi"
// @Failure 409 {object} map[string]string "Ghế đã được người khác chọn (Conflict)"
// @Failure 500 {object} map[string]string "Lỗi máy chủ nội bộ (ví dụ: lỗi định dạng ID người dùng)"
// @Router /bookings [post]
func CreateBookingController(c *gin.Context) {
	var input services.CreateBookingInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"lỗi": "Dữ liệu đầu vào không hợp lệ: " + err.Error()})
		return
	}

	if err := validate.Struct(input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"lỗi": "Lỗi xác thực dữ liệu: " + err.Error()})
		return
	}

	userID, exists := c.Get("userId")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"lỗi": "Không tìm thấy thông tin người dùng. Vui lòng đăng nhập lại."})
		return
	}
	userIDStr, ok := userID.(string)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"lỗi": "Lỗi định dạng ID người dùng từ token."})
		return
	}

	booking, err := services.CreateBooking(input, userIDStr)
	if err != nil {
		if strings.Contains(err.Error(), "không tìm thấy chuyến đi") {
			c.JSON(http.StatusNotFound, gin.H{"lỗi": err.Error()})
			return
		}
		if strings.Contains(err.Error(), "một hoặc nhiều ghế đã được người khác chọn") {
			c.JSON(http.StatusConflict, gin.H{"lỗi": err.Error()})
			return
		}
		if strings.Contains(err.Error(), "ID chuyến đi không hợp lệ") || strings.Contains(err.Error(), "ID người dùng không hợp lệ"){
			 c.JSON(http.StatusBadRequest, gin.H{"lỗi": err.Error()})
             return
		}
		c.JSON(http.StatusBadRequest, gin.H{"lỗi": "Không thể xử lý yêu cầu giữ chỗ: " + err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"thông báo": "Giữ chỗ thành công! Vui lòng tiến hành thanh toán.",
		"dữ_liệu":   booking,
	})
}

// @Summary Lấy chi tiết một booking
// @Description Lấy thông tin chi tiết của một booking dựa trên ID. Yêu cầu token xác thực và booking phải thuộc về người dùng đang đăng nhập.
// @Tags Bookings
// @Produce  json
// @Security BearerAuth
// @Param   bookingId path string true "ID của Booking" Format(ObjectID)
// @Success 200 {object} map[string]interface{} "Lấy chi tiết booking thành công. Body: {thông báo: string, dữ_liệu: models.Booking}"
// @Failure 400 {object} map[string]string "ID booking không hợp lệ"
// @Failure 401 {object} map[string]string "Yêu cầu token xác thực hoặc không thể xác định người dùng"
// @Failure 404 {object} map[string]string "Không tìm thấy booking hoặc không có quyền xem"
// @Failure 500 {object} map[string]string "Lỗi máy chủ nội bộ"
// @Router /bookings/{bookingId} [get]
func GetBookingDetailsController(c *gin.Context) {
	bookingIDStr := c.Param("bookingId")
	userID, _ := c.Get("userId")
	userIDStr, ok := userID.(string)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"lỗi": "Không thể xác định người dùng."})
		return
	}

	booking, err := services.GetBookingDetailsByID(bookingIDStr, userIDStr)
	if err != nil {
		if strings.Contains(err.Error(), "không tìm thấy booking") || strings.Contains(err.Error(), "không có quyền xem") {
			c.JSON(http.StatusNotFound, gin.H{"lỗi": err.Error()})
			return
		}
		if strings.Contains(err.Error(), "ID booking không hợp lệ") || strings.Contains(err.Error(), "ID người dùng không hợp lệ") {
			c.JSON(http.StatusBadRequest, gin.H{"lỗi": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"lỗi": "Lỗi máy chủ khi lấy chi tiết booking."})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"thông báo": "Lấy chi tiết booking thành công!",
		"dữ_liệu":   booking, 
	})
}

// @Summary Lấy lịch sử đặt vé của người dùng hiện tại
// @Description Lấy danh sách tất cả các booking của người dùng đang đăng nhập. Yêu cầu token xác thực.
// @Tags Bookings
// @Produce  json
// @Security BearerAuth
// @Success 200 {object} map[string]interface{} "Lấy lịch sử booking thành công. Body: {thông báo: string, dữ_liệu: []models.Booking}"
// @Failure 401 {object} map[string]string "Yêu cầu token xác thực hoặc không thể xác định người dùng"
// @Failure 500 {object} map[string]string "Lỗi máy chủ nội bộ"
// @Router /bookings/my [get]
func GetMyBookingsController(c *gin.Context) {
	userID, _ := c.Get("userId")
	userIDStr, ok := userID.(string)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"lỗi": "Không thể xác định người dùng."})
		return
	}

	bookings, err := services.GetBookingsByUserID(userIDStr)
	if err != nil {
		if strings.Contains(err.Error(), "ID người dùng không hợp lệ") {
             c.JSON(http.StatusUnauthorized, gin.H{"lỗi": err.Error()}) 
             return
        }
		c.JSON(http.StatusInternalServerError, gin.H{"lỗi": "Lỗi máy chủ khi lấy lịch sử booking."})
		return
	}

	if bookings == nil {
		bookings = []models.Booking{}
	}

	c.JSON(http.StatusOK, gin.H{
		"thông báo": "Lấy lịch sử booking thành công!",
		"dữ_liệu":   bookings, 
	})
}