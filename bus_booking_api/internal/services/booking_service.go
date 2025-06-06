package services

import (
	"bus_booking_api/internal/models"
	"bus_booking_api/internal/repositories"
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"math/rand"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type BookingService interface {
	CreateBooking(ctx context.Context, userIDStr string, tripIDStr string, passengers []models.PassengerInfo) (*models.Booking, error)
	ProcessMockPayment(ctx context.Context, bookingIDStr string, userIDStr string) (*models.Booking, error)
	GetBookingDetails(ctx context.Context, bookingIDStr string, userIDStr string) (*models.Booking, error) 
}

type bookingService struct {
	bookingRepo repositories.BookingRepository
	tripRepo    repositories.TripRepository
	userRepo    repositories.UserRepository 
}

func NewBookingService(bookingRepo repositories.BookingRepository, tripRepo repositories.TripRepository, userRepo repositories.UserRepository) BookingService {
	return &bookingService{
		bookingRepo: bookingRepo,
		tripRepo:    tripRepo,
		userRepo:    userRepo,
	}
}

func generateTicketCode() string {
	const charset = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	const codeLength = 8
	seededRand := rand.New(rand.NewSource(time.Now().UnixNano()))
	
	b := make([]byte, codeLength)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}

func (s *bookingService) CreateBooking(ctx context.Context, userIDStr string, tripIDStr string, passengers []models.PassengerInfo) (*models.Booking, error) {
	userID, err := primitive.ObjectIDFromHex(userIDStr)
	if err != nil {
		return nil, errors.New("ID người dùng không hợp lệ")
	}
	tripID, err := primitive.ObjectIDFromHex(tripIDStr)
	if err != nil {
		return nil, errors.New("ID chuyến đi không hợp lệ")
	}

	if len(passengers) == 0 {
		return nil, errors.New("cần ít nhất một hành khách để đặt vé")
	}

	trip, err := s.tripRepo.GetTripByID(ctx, tripID)
	if err != nil {
		log.Printf("Lỗi khi lấy thông tin chuyến đi %s: %v", tripIDStr, err)
		return nil, errors.New("lỗi khi lấy thông tin chuyến đi")
	}
	if trip == nil {
		return nil, errors.New("không tìm thấy chuyến đi")
	}
	if trip.Status != models.TripStatusScheduled {
		return nil, errors.New("chuyến đi không còn khả dụng để đặt vé")
	}

	var totalAmount float64 = 0
	seatNumbersToBook := make(map[string]bool)

	for _, p := range passengers {
		if p.SeatNumber == "" {
			return nil, errors.New("số ghế không được để trống cho hành khách")
		}
		seatNumbersToBook[p.SeatNumber] = true
		totalAmount += trip.Price 
	}

	for _, seat := range trip.Seats {
		if _, ok := seatNumbersToBook[seat.SeatNumber]; ok { 
			if seat.Status != models.SeatStatusAvailable {
				return nil, fmt.Errorf("ghế %s không còn trống hoặc đang được giữ", seat.SeatNumber)
			}
		}
	}

	newBooking := &models.Booking{
		UserID:        userID,
		TripID:        tripID,
		BookingTime:   time.Now(),
		Status:        models.BookingStatusPending, 
		PaymentStatus: models.PaymentStatusPending,
		TotalAmount:   totalAmount,
		Passengers:    passengers,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	createdBooking, err := s.bookingRepo.CreateBooking(ctx, newBooking)
	if err != nil {
		log.Printf("Lỗi khi tạo booking trong service: %v", err)
		return nil, errors.New("không thể tạo đặt vé: " + err.Error())
	}

	for _, p := range passengers {
		errUpdateSeat := s.tripRepo.UpdateTripSeatStatus(ctx, tripID, p.SeatNumber, models.SeatStatusHeld, &createdBooking.ID)
		if errUpdateSeat != nil {
			log.Printf("LỖI NGHIÊM TRỌNG: Không thể cập nhật trạng thái ghế %s cho booking %s: %v. Cần xử lý rollback!", p.SeatNumber, createdBooking.ID.Hex(), errUpdateSeat)
			return nil, errors.New("lỗi hệ thống khi giữ chỗ, vui lòng thử lại")
		}
	}
	log.Printf("Đã tạo booking %s và giữ các ghế liên quan.", createdBooking.ID.Hex())

	return createdBooking, nil
}

func (s *bookingService) ProcessMockPayment(ctx context.Context, bookingIDStr string, userIDStr string) (*models.Booking, error) {
	bookingID, err := primitive.ObjectIDFromHex(bookingIDStr)
	if err != nil {
		return nil, errors.New("ID đặt vé không hợp lệ")
	}
	userID, err := primitive.ObjectIDFromHex(userIDStr) 
	if err != nil {
		return nil, errors.New("ID người dùng không hợp lệ")
	}

	booking, err := s.bookingRepo.GetBookingByID(ctx, bookingID)
	if err != nil {
		log.Printf("Lỗi khi lấy booking %s để thanh toán: %v", bookingIDStr, err)
		return nil, errors.New("lỗi khi xử lý thanh toán")
	}
	if booking == nil {
		return nil, errors.New("không tìm thấy đặt vé để thanh toán")
	}

	if booking.UserID != userID {
		return nil, errors.New("bạn không có quyền thanh toán cho đặt vé này")
	}

	if booking.Status == models.BookingStatusConfirmed {
		return nil, errors.New("đặt vé này đã được thanh toán thành công")
	}
	if booking.Status == models.BookingStatusCancelled || booking.Status == models.BookingStatusExpired {
		return nil, errors.New("đặt vé này đã bị hủy hoặc hết hạn, không thể thanh toán")
	}

	newPaymentStatus := models.PaymentStatusPaid
	newBookingStatus := models.BookingStatusConfirmed
	mockPaymentMethod := "mock_payment_gateway"
	mockTxID := "MOCK_TX_" + bookingID.Hex()
	ticketCode := generateTicketCode() 

	err = s.bookingRepo.UpdateBookingStatus(ctx, bookingID, newBookingStatus, &newPaymentStatus, &mockPaymentMethod, &mockTxID)
	if err != nil {
		log.Printf("Lỗi khi cập nhật trạng thái booking %s sau thanh toán: %v", bookingIDStr, err)
		return nil, errors.New("lỗi khi cập nhật thông tin thanh toán")
	}

	for _, p := range booking.Passengers {
		errUpdateSeat := s.tripRepo.UpdateTripSeatStatus(ctx, booking.TripID, p.SeatNumber, models.SeatStatusBooked, &booking.ID)
		if errUpdateSeat != nil {
			log.Printf("LỖI: Không thể cập nhật trạng thái ghế %s (trip %s) thành booked sau thanh toán booking %s: %v", p.SeatNumber, booking.TripID.Hex(), booking.ID.Hex(), errUpdateSeat)
		}
	}
	log.Printf("Booking %s đã được thanh toán thành công, mã vé: %s. Các ghế đã được cập nhật.", booking.ID.Hex(), ticketCode)

	updatedBooking, _ := s.bookingRepo.GetBookingByID(ctx, bookingID)
	return updatedBooking, nil
}


func (s *bookingService) GetBookingDetails(ctx context.Context, bookingIDStr string, userIDStr string) (*models.Booking, error) {
	bookingID, err := primitive.ObjectIDFromHex(bookingIDStr)
	if err != nil {
		return nil, errors.New("ID đặt vé không hợp lệ")
	}
	userID, err := primitive.ObjectIDFromHex(userIDStr)
	if err != nil {
		return nil, errors.New("ID người dùng không hợp lệ")
	}

	booking, err := s.bookingRepo.GetBookingByID(ctx, bookingID)
	if err != nil {
		log.Printf("Lỗi khi lấy chi tiết booking %s: %v", bookingIDStr, err)
		return nil, errors.New("lỗi khi lấy thông tin đặt vé")
	}
	if booking == nil {
		return nil, errors.New("không tìm thấy đặt vé")
	}

	if booking.UserID != userID {
		return nil, errors.New("bạn không có quyền xem thông tin đặt vé này")
	}

	return booking, nil
}