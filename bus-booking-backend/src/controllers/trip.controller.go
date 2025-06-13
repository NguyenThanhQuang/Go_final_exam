package controllers

import (
	"net/http"
	"strings"

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

	var trips []models.Trip
	var err error

	if from != "" && to != "" && date != "" {
		trips, err = services.SearchTrips(from, to, date)
	} else if from == "" && to == "" && date == "" {
		trips, err = services.GetAllTrips() 
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"lỗi": "Cần cung cấp đủ các tham số 'from', 'to', 'date' cho việc tìm kiếm, hoặc không cung cấp tham số nào để lấy tất cả chuyến đi."})
		return
	}


	if err != nil {
		if strings.Contains(err.Error(), "Lỗi phân tích ngày") {
             c.JSON(http.StatusBadRequest, gin.H{"lỗi": "Định dạng ngày không hợp lệ. Vui lòng sử dụng YYYY-MM-DD."})
             return
        }
		c.JSON(http.StatusInternalServerError, gin.H{"lỗi": "Đã có lỗi xảy ra ở máy chủ khi xử lý yêu cầu chuyến đi."})
		return
	}

	if len(trips) == 0 {
		c.JSON(http.StatusOK, gin.H{
			"thông báo": "Không tìm thấy chuyến đi nào phù hợp.",
			"dữ_liệu":   []models.Trip{},
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"thông báo": "Yêu cầu chuyến đi thành công!", 
		"dữ_liệu": trips,
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
		"dữ_liệu": trip,
	})
}

