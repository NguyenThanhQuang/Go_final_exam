package services

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/Go_final_exam/bus-booking-backend/src/config"
	"github.com/Go_final_exam/bus-booking-backend/src/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func SearchTrips(from, to, date string) ([]models.Trip, error) {
	tripCollection := config.DB.Collection("trips")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := bson.M{}

	if from != "" {
		filter["route.from.name"] = from
	}
	if to != "" {
		filter["route.to.name"] = to
	}

	if date != "" {
		layout := "2006-01-02" 
		startOfDay, err := time.Parse(layout, date)
		if err != nil {
			log.Printf("Lỗi phân tích ngày: %v", err)
			return nil, err
		}
		endOfDay := startOfDay.Add(24 * time.Hour)

		filter["departureTime"] = bson.M{
			"$gte": startOfDay,
			"$lt":  endOfDay,
		}
	}

	findOptions := options.Find()
	findOptions.SetSort(bson.D{{Key: "departureTime", Value: 1}})

	cursor, err := tripCollection.Find(ctx, filter, findOptions)
	if err != nil {
		log.Printf("Lỗi khi tìm kiếm chuyến đi: %v", err)
		return nil, err
	}
	defer cursor.Close(ctx)

	var trips []models.Trip
	if err = cursor.All(ctx, &trips); err != nil {
		log.Printf("Lỗi khi đọc dữ liệu chuyến đi từ cursor: %v", err)
		return nil, err
	}

	for i := range trips {
		availableCount := 0
		for _, seat := range trips[i].Seats {
			if seat.Status == "available" {
				availableCount++
			}
		}
		trips[i].AvailableSeats = availableCount
	}

	return trips, nil
}

func GetTripByID(tripID string) (*models.Trip, error) {
	tripCollection := config.DB.Collection("trips")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	objID, err := primitive.ObjectIDFromHex(tripID)
	if err != nil {
		log.Printf("ID chuyến đi không hợp lệ: %s", tripID)
		return nil, errors.New("ID chuyến đi không hợp lệ")
	}

	var trip models.Trip
	err = tripCollection.FindOne(ctx, bson.M{"_id": objID}).Decode(&trip)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			log.Printf("Không tìm thấy chuyến đi với ID: %s", tripID)
			return nil, errors.New("không tìm thấy chuyến đi")
		}
		log.Printf("Lỗi khi tìm chuyến đi bằng ID: %v", err)
		return nil, errors.New("lỗi máy chủ khi truy vấn dữ liệu")
	}

	availableCount := 0
	for _, seat := range trip.Seats {
		if seat.Status == "available" {
			availableCount++
		}
	}
	trip.AvailableSeats = availableCount

	return &trip, nil
}

func GetAllTrips() ([]models.Trip, error) {
	tripCollection := config.DB.Collection("trips")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	findOptions := options.Find()
	findOptions.SetSort(bson.D{{Key: "departureTime", Value: 1}}) 

	cursor, err := tripCollection.Find(ctx, bson.M{}, findOptions) 
	if err != nil {
		log.Printf("Lỗi khi lấy tất cả chuyến đi: %v", err)
		return nil, errors.New("lỗi máy chủ khi truy vấn tất cả chuyến đi")
	}
	defer cursor.Close(ctx)

	var trips []models.Trip
	if err = cursor.All(ctx, &trips); err != nil {
		log.Printf("Lỗi khi đọc dữ liệu tất cả chuyến đi từ cursor: %v", err)
		return nil, errors.New("lỗi máy chủ khi đọc dữ liệu chuyến đi")
	}

	for i := range trips {
		availableCount := 0
		for _, seat := range trips[i].Seats {
			if seat.Status == "available" {
				availableCount++
			}
		}
		trips[i].AvailableSeats = availableCount
	}

	return trips, nil
}
