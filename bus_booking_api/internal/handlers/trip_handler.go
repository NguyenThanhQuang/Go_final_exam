package handlers

import (
	"bus_booking_api/internal/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

type TripHandler struct {
	tripService services.TripService
}

func NewTripHandler(tripService services.TripService) *TripHandler {
	return &TripHandler{
		tripService: tripService,
	}
}

// SearchTrips godoc
// @Summary Tìm kiếm chuyến đi
// @Description Tìm kiếm các chuyến đi dựa trên điểm đi, điểm đến và ngày.
// @Tags trips
// @Accept  json
// @Produce  json
// @Param from query string true "Tên điểm đi (ví dụ: Sài Gòn)"
// @Param to query string true "Tên điểm đến (ví dụ: Đà Lạt)"
// @Param date query string true "Ngày khởi hành (YYYY-MM-DD)"
// @Success 200 {array} models.Trip "Danh sách chuyến đi phù hợp"
// @Failure 400 {object} map[string]string "Tham số không hợp lệ"
// @Failure 500 {object} map[string]string "Lỗi máy chủ nội bộ"
// @Router /trips/search [get]
func (h *TripHandler) SearchTrips(c *gin.Context) {
	from := c.Query("from")
	to := c.Query("to")
	dateStr := c.Query("date") // Format: YYYY-MM-DD

	if from == "" || to == "" || dateStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Vui lòng cung cấp điểm đi, điểm đến và ngày tìm kiếm."})
		return
	}

	trips, err := h.tripService.SearchTrips(c.Request.Context(), from, to, dateStr)
	if err != nil {
		if err.Error() == "định dạng ngày không hợp lệ, vui lòng sử dụng YYYY-MM-DD" ||
			err.Error() == "điểm đi, điểm đến và ngày không được để trống" {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Không thể tìm kiếm chuyến đi: " + err.Error()})
		}
		return
	}

	if len(trips) == 0 {
		c.JSON(http.StatusOK, gin.H{"message": "Không tìm thấy chuyến đi nào phù hợp.", "data": []interface{}{}})
		return
	}

	c.JSON(http.StatusOK, trips)
}

// GetTripDetails godoc
// @Summary Lấy chi tiết một chuyến đi
// @Description Lấy thông tin chi tiết của một chuyến đi dựa vào ID.
// @Tags trips
// @Accept  json
// @Produce  json
// @Param   tripId path string true "ID của chuyến đi"
// @Success 200 {object} models.Trip "Chi tiết chuyến đi"
// @Failure 400 {object} map[string]string "ID chuyến đi không hợp lệ"
// @Failure 404 {object} map[string]string "Không tìm thấy chuyến đi"
// @Failure 500 {object} map[string]string "Lỗi máy chủ nội bộ"
// @Router /trips/{tripId} [get]
func (h *TripHandler) GetTripDetails(c *gin.Context) {
	tripID := c.Param("tripId")

	trip, err := h.tripService.GetTripDetails(c.Request.Context(), tripID)
	if err != nil {
		if err.Error() == "ID chuyến đi không hợp lệ" {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		} else if err.Error() == "không tìm thấy chuyến đi" {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Không thể lấy thông tin chuyến đi: " + err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, trip)
}