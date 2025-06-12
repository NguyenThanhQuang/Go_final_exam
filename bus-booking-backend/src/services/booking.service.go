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
