package services

import (
	"bus_booking_api/internal/models"
	"bus_booking_api/internal/repositories"
	"context"
	"errors"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type TripService interface {
	SearchTrips(ctx context.Context, from, to, dateStr string) ([]models.Trip, error)
	GetTripDetails(ctx context.Context, tripIDStr string) (*models.Trip, error)
	AddTrip(ctx context.Context, tripData models.Trip) (*models.Trip, error) 
}

type tripService struct {
	tripRepo    repositories.TripRepository
	companyRepo repositories.CompanyRepository 
	vehicleRepo repositories.VehicleRepository 
}

func NewTripService(tripRepo repositories.TripRepository, companyRepo repositories.CompanyRepository, vehicleRepo repositories.VehicleRepository) TripService {
	return &tripService{
		tripRepo:    tripRepo,
		companyRepo: companyRepo,
		vehicleRepo: vehicleRepo,
	}
}

func (s *tripService) SearchTrips(ctx context.Context, from, to, dateStr string) ([]models.Trip, error) {
	if from == "" || to == "" || dateStr == "" {
		return nil, errors.New("điểm đi, điểm đến và ngày không được để trống")
	}

	parsedDate, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		return nil, errors.New("định dạng ngày không hợp lệ, vui lòng sử dụng YYYY-MM-DD")
	}

	trips, err := s.tripRepo.SearchTrips(ctx, from, to, parsedDate)
	if err != nil {
		log.Printf("Lỗi khi tìm kiếm chuyến đi từ repository: %v", err)
		return nil, errors.New("đã xảy ra lỗi khi tìm kiếm chuyến đi")
	}
	return trips, nil
}

func (s *tripService) GetTripDetails(ctx context.Context, tripIDStr string) (*models.Trip, error) {
	tripID, err := primitive.ObjectIDFromHex(tripIDStr)
	if err != nil {
		return nil, errors.New("ID chuyến đi không hợp lệ")
	}

	trip, err := s.tripRepo.GetTripByID(ctx, tripID)
	if err != nil {
		log.Printf("Lỗi khi lấy chi tiết chuyến đi từ repository: %v", err)
		return nil, errors.New("đã xảy ra lỗi khi lấy thông tin chuyến đi")
	}
	if trip == nil {
		return nil, errors.New("không tìm thấy chuyến đi")
	}

	return trip, nil
}

func (s *tripService) AddTrip(ctx context.Context, tripData models.Trip) (*models.Trip, error) {
    createdTrip, err := s.tripRepo.CreateTrip(ctx, &tripData)
    if err != nil {
        log.Printf("Lỗi khi tạo chuyến đi trong service: %v", err)
        return nil, errors.New("không thể tạo chuyến đi: " + err.Error())
    }
    return createdTrip, nil
}