package controllers

import (
	"net/http"

	"github.com/Go_final_exam/bus-booking-backend/src/models"
	"github.com/Go_final_exam/bus-booking-backend/src/services"
	"github.com/gin-gonic/gin"
)

// @Summary Tìm kiếm chuyến đi
// @Description Tìm kiếm các chuyến đi dựa trên điểm đi, điểm đến và ngày khởi hành.
// @Tags Trips
// @Accept  json
// @Produce  json
// @Param from query string true "Tên điểm đi (Ví dụ: 'TP. Hồ Chí Minh')"
// @Param to query string true "Tên điểm đến (Ví dụ: 'Đà Lạt')"
// @Param date query string true "Ngày đi theo định dạng YYYY-MM-DD (Ví dụ: '2024-05-25')"
// @Success 200 {object} map[string]interface{} "Danh sách các chuyến đi phù hợp"
// @Failure 400 {object} map[string]string "Các tham số query bắt buộc bị thiếu"
// @Failure 500 {object} map[string]string "Lỗi máy chủ"
// @Router /trips [get]
func SearchTripsController(c *gin.Context) {
	from := c.Query("from")
	to := c.Query("to")
	date := c.Query("date")

	if from == "" || to == "" || date == "" {
		c.JSON(http.StatusBadRequest, gin.H{"lỗi": "Các tham số 'from', 'to', và 'date' là bắt buộc."})
		return
	}

	trips, err := services.SearchTrips(from, to, date)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"lỗi": "Đã có lỗi xảy ra ở máy chủ khi tìm kiếm chuyến đi."})
		return
	}

	if len(trips) == 0 {
		c.JSON(http.StatusOK, gin.H{
			"thông báo": "Không tìm thấy chuyến đi nào phù hợp với yêu cầu của bạn.",
			"dữ liệu":   []models.Trip{},
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"thông báo": "Tìm kiếm chuyến đi thành công!",
		"dữ liệu":   trips,
	})
}

// @Summary Lấy thông tin chi tiết một chuyến đi
// @Description Lấy toàn bộ thông tin chi tiết của một chuyến đi, bao gồm cả sơ đồ ghế.
// @Tags Trips
// @Accept  json
// @Produce  json
// @Param tripId path string true "ID của chuyến đi"
// @Success 200 {object} map[string]interface{} "Thông tin chi tiết của chuyến đi"
// @Failure 400 {object} map[string]string "ID chuyến đi không hợp lệ"
// @Failure 404 {object} map[string]string "Không tìm thấy chuyến đi"
// @Router /trips/{tripId} [get]
func GetTripDetailsController(c *gin.Context) {
	tripID := c.Param("tripId")

	trip, err := services.GetTripByID(tripID)
	if err != nil {
		if err.Error() == "không tìm thấy chuyến đi" {
			c.JSON(http.StatusNotFound, gin.H{"lỗi": err.Error()})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"lỗi": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"thông báo": "Lấy thông tin chi tiết chuyến đi thành công!",
		"dữ liệu":   trip,
	})
}
