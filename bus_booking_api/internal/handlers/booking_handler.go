package handlers

import (
	"bus_booking_api/internal/models"
	"bus_booking_api/internal/services"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type BookingHandler struct {
	bookingService services.BookingService
}

func NewBookingHandler(bookingService services.BookingService) *BookingHandler {
	return &BookingHandler{
		bookingService: bookingService,
	}
}

// CreateBooking godoc
// @Summary Tạo một đặt vé mới
// @Description Người dùng đặt vé cho một chuyến đi với thông tin hành khách.
// @Tags bookings
// @Security ApiKeyAuth
// @Accept  json
// @Produce  json
// @Param   bookingRequest body models.CreateBookingRequest true "Thông tin đặt vé"
// @Success 201 {object} models.Booking "Đặt vé thành công, chờ thanh toán"
// @Failure 400 {object} map[string]string "Dữ liệu không hợp lệ"
// @Failure 401 {object} map[string]string "Chưa xác thực hoặc token không hợp lệ"
// @Failure 404 {object} map[string]string "Chuyến đi không tìm thấy"
// @Failure 409 {object} map[string]string "Ghế không còn trống"
// @Failure 500 {object} map[string]string "Lỗi máy chủ nội bộ"
// @Router /bookings [post]
func (h *BookingHandler) CreateBooking(c *gin.Context) {
	var req models.CreateBookingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Dữ liệu không hợp lệ: " + err.Error()})
		return
	}

	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Người dùng chưa được xác thực."})
		return
	}
	userIDStr, ok := userID.(string)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Lỗi xử lý thông tin người dùng."})
		return
	}

	booking, err := h.bookingService.CreateBooking(c.Request.Context(), userIDStr, req.TripID, req.Passengers)
	if err != nil {
		errMsg := err.Error()
		if errMsg == "ID người dùng không hợp lệ" || errMsg == "ID chuyến đi không hợp lệ" ||
			errMsg == "cần ít nhất một hành khách để đặt vé" || errMsg == "số ghế không được để trống cho hành khách" {
			c.JSON(http.StatusBadRequest, gin.H{"error": errMsg})
		} else if errMsg == "không tìm thấy chuyến đi" {
			c.JSON(http.StatusNotFound, gin.H{"error": errMsg})
		} else if errMsg == "chuyến đi không còn khả dụng để đặt vé" || strings.Contains(errMsg, "không còn trống hoặc đang được giữ") {
			c.JSON(http.StatusConflict, gin.H{"error": errMsg}) // 409 Conflict
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Không thể tạo đặt vé: " + errMsg})
		}
		return
	}

	c.JSON(http.StatusCreated, booking)
}

// MockPayment godoc
// @Summary Giả lập thanh toán cho một đặt vé
// @Description Đánh dấu một đặt vé là đã thanh toán thành công (mock).
// @Tags bookings
// @Security ApiKeyAuth
// @Accept  json
// @Produce  json
// @Param   paymentRequest body models.MockPaymentRequest true "ID của đặt vé cần thanh toán"
// @Success 200 {object} models.Booking "Thanh toán thành công, đặt vé đã được xác nhận"
// @Failure 400 {object} map[string]string "Dữ liệu không hợp lệ hoặc đặt vé không phù hợp để thanh toán"
// @Failure 401 {object} map[string]string "Chưa xác thực hoặc không có quyền"
// @Failure 404 {object} map[string]string "Đặt vé không tìm thấy"
// @Failure 500 {object} map[string]string "Lỗi máy chủ nội bộ"
// @Router /bookings/payment/mock [post]
func (h *BookingHandler) MockPayment(c *gin.Context) {
	var req models.MockPaymentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Dữ liệu không hợp lệ: " + err.Error()})
		return
	}

	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Người dùng chưa được xác thực."})
		return
	}
	userIDStr, _ := userID.(string)

	booking, err := h.bookingService.ProcessMockPayment(c.Request.Context(), req.BookingID, userIDStr)
	if err != nil {
		errMsg := err.Error()
		if errMsg == "ID đặt vé không hợp lệ" || errMsg == "ID người dùng không hợp lệ" ||
			errMsg == "đặt vé này đã được thanh toán thành công" ||
			errMsg == "đặt vé này đã bị hủy hoặc hết hạn, không thể thanh toán" {
			c.JSON(http.StatusBadRequest, gin.H{"error": errMsg})
		} else if errMsg == "không tìm thấy đặt vé để thanh toán" {
			c.JSON(http.StatusNotFound, gin.H{"error": errMsg})
		} else if errMsg == "bạn không có quyền thanh toán cho đặt vé này" {
			c.JSON(http.StatusForbidden, gin.H{"error": errMsg}) // 403 Forbidden
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Lỗi khi xử lý thanh toán: " + errMsg})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Thanh toán thành công!", "booking": booking})
}

// GetBookingDetails godoc
// @Summary Lấy chi tiết một đặt vé
// @Description Người dùng lấy thông tin chi tiết của một đặt vé họ đã tạo.
// @Tags bookings
// @Security ApiKeyAuth
// @Accept  json
// @Produce  json
// @Param   bookingId path string true "ID của đặt vé"
// @Success 200 {object} models.Booking "Chi tiết đặt vé"
// @Failure 400 {object} map[string]string "ID đặt vé không hợp lệ"
// @Failure 401 {object} map[string]string "Chưa xác thực"
// @Failure 403 {object} map[string]string "Không có quyền xem đặt vé này"
// @Failure 404 {object} map[string]string "Đặt vé không tìm thấy"
// @Failure 500 {object} map[string]string "Lỗi máy chủ nội bộ"
// @Router /bookings/{bookingId} [get]
func (h *BookingHandler) GetBookingDetails(c *gin.Context) {
	bookingID := c.Param("bookingId")

	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Người dùng chưa được xác thực."})
		return
	}
	userIDStr, _ := userID.(string)

	booking, err := h.bookingService.GetBookingDetails(c.Request.Context(), bookingID, userIDStr)
	if err != nil {
		errMsg := err.Error()
		if errMsg == "ID đặt vé không hợp lệ" || errMsg == "ID người dùng không hợp lệ" {
			c.JSON(http.StatusBadRequest, gin.H{"error": errMsg})
		} else if errMsg == "không tìm thấy đặt vé" {
			c.JSON(http.StatusNotFound, gin.H{"error": errMsg})
		} else if errMsg == "bạn không có quyền xem thông tin đặt vé này" {
			c.JSON(http.StatusForbidden, gin.H{"error": errMsg})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Không thể lấy thông tin đặt vé: " + errMsg})
		}
		return
	}

	c.JSON(http.StatusOK, booking)
}