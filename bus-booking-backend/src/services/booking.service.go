package services

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/Go_final_exam/bus-booking-backend/src/config"
	"github.com/Go_final_exam/bus-booking-backend/src/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type CreateBookingInput struct {
	TripID      string   `json:"tripId" validate:"required"`
	SeatNumbers []string `json:"seatNumbers" validate:"required,min=1"`
}

func CreateBooking(input CreateBookingInput, userIDStr string) (*models.Booking, error) {
    ctx := context.Background()
    tripCollection := config.DB.Collection("trips")
    bookingCollection := config.DB.Collection("bookings")

    tripID, err := primitive.ObjectIDFromHex(input.TripID)
    if err != nil {
        return nil, errors.New("ID chuyến đi không hợp lệ")
    }
    userID, err := primitive.ObjectIDFromHex(userIDStr)
    if err != nil {
        return nil, errors.New("ID người dùng không hợp lệ")
    }

    var trip models.Trip
    err = tripCollection.FindOne(ctx, bson.M{"_id": tripID}).Decode(&trip)
    if err != nil {
        if err == mongo.ErrNoDocuments {
            return nil, errors.New("không tìm thấy chuyến đi")
        }
        log.Printf("Lỗi khi FindOne trip: %v", err)
        return nil, errors.New("lỗi khi tìm kiếm chuyến đi")
    }

    seatPrice := trip.Price
    seatsToBook := make(map[string]bool)
    for _, s := range input.SeatNumbers {
        seatsToBook[s] = true
    }

    for _, seat := range trip.Seats {
        if seatsToBook[seat.SeatNumber] {
            if seat.Status != "available" {
                return nil, fmt.Errorf("ghế '%s' đã được chọn hoặc không còn trống", seat.SeatNumber)
            }
        }
    }

    newBooking := models.Booking{
        ID:          primitive.NewObjectID(),
        UserID:      userID,
        TripID:      tripID,
        BookingTime: time.Now(),
        Status:      "held",
        TotalAmount: float64(len(input.SeatNumbers)) * seatPrice,
        Passengers:  []models.Passenger{},
        CreatedAt:   time.Now(),
        UpdatedAt:   time.Now(),
    }
    for _, seatNum := range input.SeatNumbers {
        newBooking.Passengers = append(newBooking.Passengers, models.Passenger{SeatNumber: seatNum})
    }

    _, err = bookingCollection.InsertOne(ctx, newBooking)
    if err != nil {
        return nil, errors.New("không thể tạo booking mới")
    }

    filter := bson.M{
        "_id":              tripID,
        "seats.seatNumber": bson.M{"$in": input.SeatNumbers},
    }
    update := bson.M{
        "$set": bson.M{"seats.$[elem].status": "held"},
    }
    arrayFilters := options.ArrayFilters{
        Filters: []interface{}{bson.M{"elem.seatNumber": bson.M{"$in": input.SeatNumbers}}},
    }
    updateOptions := options.UpdateOptions{ArrayFilters: &arrayFilters}

    _, err = tripCollection.UpdateOne(ctx, filter, update, &updateOptions)
    if err != nil {
        return nil, errors.New("lỗi khi cập nhật trạng thái ghế")
    }

    return &newBooking, nil
}

func GetBookingDetailsByID(bookingIDStr string, userIDStr string) (*models.Booking, error) {
	bookingCollection := config.DB.Collection("bookings")
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second) // Tăng timeout một chút cho aggregation
	defer cancel()

	log.Printf("BACKEND DEBUG: GetBookingDetailsByID - Bắt đầu cho bookingID: %s, userID: %s", bookingIDStr, userIDStr)

	objBookingID, err := primitive.ObjectIDFromHex(bookingIDStr)
	if err != nil {
		log.Printf("BACKEND ERROR: GetBookingDetailsByID - ID booking không hợp lệ: %s, lỗi: %v", bookingIDStr, err)
		return nil, errors.New("ID booking không hợp lệ")
	}
	objUserID, err := primitive.ObjectIDFromHex(userIDStr)
	if err != nil {
		// Lỗi này không nên xảy ra nếu userID từ JWT là hợp lệ
		log.Printf("BACKEND ERROR: GetBookingDetailsByID - Lỗi chuyển đổi userID từ JWT: %s, lỗi: %v", userIDStr, err)
		return nil, errors.New("ID người dùng không hợp lệ")
	}

	// Sử dụng Aggregation Pipeline để lookup thông tin Trip
	pipeline := mongo.Pipeline{
		{{"$match", bson.D{{"_id", objBookingID}, {"userId", objUserID}}}}, // Khớp bookingId và userId
		{{"$lookup", bson.D{
			{"from", "trips"},         // Tên collection trips (đảm bảo đúng tên)
			{"localField", "tripId"},  // Trường trong collection bookings (đảm bảo đúng tên và kiểu)
			{"foreignField", "_id"},   // Trường trong collection trips (đảm bảo đúng tên và kiểu)
			{"as", "tripInfoArr"},     // Tên mảng kết quả của lookup
		}}},
		// Bước $unwind: chuyển mảng tripInfoArr (thường chỉ có 1 phần tử nếu tìm thấy) thành object
		// preserveNullAndEmptyArrays: true sẽ giữ lại document booking ngay cả khi không có trip nào khớp (tripInfoArr sẽ là mảng rỗng)
		{{"$unwind", bson.D{
			{"path", "$tripInfoArr"},
			{"preserveNullAndEmptyArrays", true},
		}}},
		// Bước $addFields: Tạo trường mới "tripInfo" từ trường "tripInfoArr" (sau khi unwind, tripInfoArr là object hoặc null)
		// Nếu tripInfoArr là null (do preserveNullAndEmptyArrays và không tìm thấy trip), tripInfo cũng sẽ là null
		{{"$addFields", bson.D{
			{"tripInfo", "$tripInfoArr"},
		}}},
		// Bước $project: Loại bỏ trường "tripInfoArr" không cần thiết nữa
		{{"$project", bson.D{
			{"tripInfoArr", 0},
		}}},
	}

	log.Printf("BACKEND DEBUG: GetBookingDetailsByID - Pipeline được xây dựng: %+v", pipeline)

	cursor, err := bookingCollection.Aggregate(ctx, pipeline)
	if err != nil {
		log.Printf("BACKEND ERROR: GetBookingDetailsByID - Lỗi thực thi aggregation: %v", err)
		return nil, errors.New("lỗi hệ thống khi truy vấn chi tiết booking")
	}
	defer cursor.Close(ctx)

	var results []models.Booking // Kết quả aggregate luôn là một mảng
	if err = cursor.All(ctx, &results); err != nil {
		log.Printf("BACKEND ERROR: GetBookingDetailsByID - Lỗi đọc tất cả kết quả từ cursor aggregation: %v", err)
		return nil, errors.New("lỗi hệ thống khi đọc dữ liệu booking")
	}

	log.Printf("BACKEND DEBUG: GetBookingDetailsByID - Số lượng results từ aggregation: %d", len(results))

	if len(results) == 0 {
		log.Printf("BACKEND DEBUG: GetBookingDetailsByID - Không tìm thấy booking nào sau aggregation (có thể do bookingId sai, userId không khớp, hoặc tripId không tìm thấy và preserveNullAndEmptyArrays là false ở bước nào đó - nhưng hiện tại là true).")
		// Kiểm tra xem booking có tồn tại nhưng không thuộc user này không (nếu $match không lọc được)
		// Tuy nhiên, $match đã nên xử lý việc này.
		// var tempBooking models.Booking
		// errCheck := bookingCollection.FindOne(ctx, bson.M{"_id": objBookingID}).Decode(&tempBooking)
		// if errCheck == nil { // Booking tồn tại
		// 	log.Printf("BACKEND WARNING: GetBookingDetailsByID - Booking %s tồn tại nhưng không thuộc user %s hoặc không qua được pipeline.", bookingIDStr, userIDStr)
		// 	return nil, errors.New("bạn không có quyền xem booking này hoặc booking không hợp lệ")
		// }
		return nil, errors.New("không tìm thấy booking hoặc bạn không có quyền xem")
	}

	// Log chi tiết từng document trong results (nếu có nhiều hơn 1, dù $match theo _id thường chỉ ra 1)
	// for i, res := range results {
	//     log.Printf("BACKEND DEBUG: GetBookingDetailsByID - Result[%d]: %+v", i, res)
	//     if res.TripInfo != nil {
	//          log.Printf("BACKEND DEBUG: GetBookingDetailsByID - Result[%d].TripInfo: %+v", i, *res.TripInfo)
	//     } else {
	//          log.Printf("BACKEND DEBUG: GetBookingDetailsByID - Result[%d].TripInfo LÀ NIL", i)
	//     }
	// }

	log.Printf("BACKEND DEBUG: GetBookingDetailsByID - Booking sẽ trả về cho client: %+v", results[0])
	if results[0].TripInfo != nil {
		log.Printf("BACKEND DEBUG: GetBookingDetailsByID - TripInfo trong booking sẽ trả về: %+v", *results[0].TripInfo)
        log.Printf("BACKEND DEBUG: GetBookingDetailsByID - TripInfo.CompanyName: %s", results[0].TripInfo)
        log.Printf("BACKEND DEBUG: GetBookingDetailsByID - TripInfo.Route.From.Name: %s", results[0].TripInfo.Route.From.Name)
	} else {
		log.Printf("BACKEND DEBUG: GetBookingDetailsByID - TripInfo trong booking sẽ trả về LÀ NIL")
	}

	return &results[0], nil
}
  func GetBookingsByUserID(userIDStr string) ([]models.Booking, error) {
  	bookingCollection := config.DB.Collection("bookings")
  	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
  	defer cancel()

  	objUserID, err := primitive.ObjectIDFromHex(userIDStr)
  	if err != nil {
  		log.Printf("Lỗi chuyển đổi userID từ JWT (GetBookingsByUserID): %v", err)
  		return nil, errors.New("ID người dùng không hợp lệ")
  	}

  	pipeline := mongo.Pipeline{
  		{{"$match", bson.D{{"userId", objUserID}}}},
        {{"$sort", bson.D{{"bookingTime", -1}}}}, 
  		{{"$lookup", bson.D{
  			{"from", "trips"},
  			{"localField", "tripId"},
  			{"foreignField", "_id"},
  			{"as", "tripInfoArr"},
  		}}},
  		{{"$unwind", bson.D{
  			{"path", "$tripInfoArr"},
  			{"preserveNullAndEmptyArrays", true},
  		}}},
        {{"$addFields", bson.D{
            {"tripInfo", "$tripInfoArr"},
        }}},
        {{"$project", bson.D{
            {"tripInfoArr", 0},
        }}},
  	}

  	cursor, err := bookingCollection.Aggregate(ctx, pipeline, options.Aggregate()) 
  	if err != nil {
  		log.Printf("Lỗi aggregation khi lấy danh sách booking: %v", err)
  		return nil, errors.New("lỗi hệ thống khi truy vấn danh sách booking")
  	}
  	defer cursor.Close(ctx)

  	var bookings []models.Booking
  	if err = cursor.All(ctx, &bookings); err != nil {
  		log.Printf("Lỗi đọc cursor aggregation (danh sách booking): %v", err)
  		return nil, errors.New("lỗi hệ thống khi đọc dữ liệu booking")
  	}

  	return bookings, nil
  }