package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

const (
	RoleUser         = "user"
	RoleCompanyAdmin = "company_admin"
	RoleAdmin        = "admin"
)

const (
	TripStatusScheduled = "scheduled"
	TripStatusDeparted  = "departed"
	TripStatusArrived   = "arrived"
	TripStatusCancelled = "cancelled"
)

const (
	SeatStatusAvailable = "available"
	SeatStatusHeld      = "held"
	SeatStatusBooked    = "booked"
)

const (
	BookingStatusPending   = "pending"
	BookingStatusHeld      = "held"
	BookingStatusConfirmed = "confirmed"
	BookingStatusCancelled = "cancelled"
	BookingStatusExpired   = "expired"
)

const (
	PaymentStatusPending = "pending"
	PaymentStatusPaid    = "paid"
	PaymentStatusFailed  = "failed"
)

type User struct {
	ID           primitive.ObjectID  `bson:"_id,omitempty" json:"id,omitempty"`
	Email        string              `bson:"email" json:"email"`
	Phone        string              `bson:"phone" json:"phone"`
	PasswordHash string              `bson:"passwordHash" json:"-"` 
	Name         string              `bson:"name" json:"name"`
	Role         string              `bson:"role" json:"role"`                                  
	CompanyID    *primitive.ObjectID `bson:"companyId,omitempty" json:"companyId,omitempty"` 
	CreatedAt    time.Time           `bson:"createdAt" json:"createdAt"`
	UpdatedAt    time.Time           `bson:"updatedAt" json:"updatedAt"`
}

type Company struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Name        string             `bson:"name" json:"name"`
	Code        string             `bson:"code" json:"code"` 
	Address     string             `bson:"address,omitempty" json:"address,omitempty"`
	Phone       string             `bson:"phone,omitempty" json:"phone,omitempty"`
	Email       string             `bson:"email,omitempty" json:"email,omitempty"`
	Description string             `bson:"description,omitempty" json:"description,omitempty"`
	LogoURL     string             `bson:"logoUrl,omitempty" json:"logoUrl,omitempty"`
	IsActive    bool               `bson:"isActive" json:"isActive"`
	CreatedAt   time.Time          `bson:"createdAt" json:"createdAt"`
	UpdatedAt   time.Time          `bson:"updatedAt" json:"updatedAt"`
}

type VehicleSeatMap struct {
	Rows    int             `bson:"rows" json:"rows"`
	Cols    int             `bson:"cols" json:"cols"`
	Layout  [][]interface{} `bson:"layout" json:"layout"` 
	Legend  map[string]string `bson:"legend,omitempty" json:"legend,omitempty"` 
}

type Vehicle struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	CompanyID   primitive.ObjectID `bson:"companyId" json:"companyId"`
	Type        string             `bson:"type" json:"type"`
	Description string             `bson:"description,omitempty" json:"description,omitempty"`
	SeatMap     VehicleSeatMap     `bson:"seatMap" json:"seatMap"`
	TotalSeats  int                `bson:"totalSeats" json:"totalSeats"`
	CreatedAt   time.Time          `bson:"createdAt" json:"createdAt"`
	UpdatedAt   time.Time          `bson:"updatedAt" json:"updatedAt"`
}

type GeoJSONPoint struct {
	Type        string    `bson:"type" json:"type"` 
	Coordinates []float64 `bson:"coordinates" json:"coordinates"`
}

type StopPoint struct {
	Name     string       `bson:"name" json:"name"`
	Location GeoJSONPoint `bson:"location" json:"location"`
}

type TripStop struct {
	Name                  string       `bson:"name" json:"name"`
	Location              GeoJSONPoint `bson:"location" json:"location"`
	ExpectedArrivalTime   *time.Time   `bson:"expectedArrivalTime,omitempty" json:"expectedArrivalTime,omitempty"`
	ExpectedDepartureTime *time.Time   `bson:"expectedDepartureTime,omitempty" json:"expectedDepartureTime,omitempty"`
}

type Route struct {
	From     StopPoint  `bson:"from" json:"from"`
	To       StopPoint  `bson:"to" json:"to"`
	Stops    []TripStop `bson:"stops,omitempty" json:"stops,omitempty"`
	Polyline string     `bson:"polyline,omitempty" json:"polyline,omitempty"` 
}

type TripSeat struct {
	SeatNumber string              `bson:"seatNumber" json:"seatNumber"`
	Status     string              `bson:"status" json:"status"` 
	BookingID  *primitive.ObjectID `bson:"bookingId,omitempty" json:"bookingId,omitempty"`
}

type Trip struct {
	ID                  primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	CompanyID           primitive.ObjectID `bson:"companyId" json:"companyId"`
	VehicleID           primitive.ObjectID `bson:"vehicleId" json:"vehicleId"`
	Route               Route              `bson:"route" json:"route"`
	DepartureTime       time.Time          `bson:"departureTime" json:"departureTime"`
	ExpectedArrivalTime time.Time          `bson:"expectedArrivalTime" json:"expectedArrivalTime"`
	Price               float64            `bson:"price" json:"price"`
	Status              string             `bson:"status" json:"status"` 
	Seats               []TripSeat         `bson:"seats" json:"seats"`
	CreatedAt           time.Time          `bson:"createdAt" json:"createdAt"`
	UpdatedAt           time.Time          `bson:"updatedAt" json:"updatedAt"`
}

type PassengerInfo struct {
	Name       string `bson:"name" json:"name"`
	Phone      string `bson:"phone" json:"phone"`
	SeatNumber string `bson:"seatNumber" json:"seatNumber"`
}

type Booking struct {
	ID                          primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	UserID                      primitive.ObjectID `bson:"userId" json:"userId"`
	TripID                      primitive.ObjectID `bson:"tripId" json:"tripId"`
	BookingTime                 time.Time          `bson:"bookingTime" json:"bookingTime"`
	Status                      string             `bson:"status" json:"status"`           
	HeldUntil                   *time.Time         `bson:"heldUntil,omitempty" json:"heldUntil,omitempty"` 
	PaymentStatus               string             `bson:"paymentStatus" json:"paymentStatus"`            
	PaymentMethod               *string            `bson:"paymentMethod,omitempty" json:"paymentMethod,omitempty"`
	TotalAmount                 float64            `bson:"totalAmount" json:"totalAmount"`
	Passengers                  []PassengerInfo    `bson:"passengers" json:"passengers"`
	TicketCode                  *string            `bson:"ticketCode,omitempty" json:"ticketCode,omitempty"` 
	PaymentGatewayTransactionID *string            `bson:"paymentGatewayTransactionId,omitempty" json:"paymentGatewayTransactionId,omitempty"` 
	CreatedAt                   time.Time          `bson:"createdAt" json:"createdAt"`
	UpdatedAt                   time.Time          `bson:"updatedAt" json:"updatedAt"`
}

type RegisterRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Phone    string `json:"phone" binding:"required"` 
	Password string `json:"password" binding:"required,min=6"`
	Name     string `json:"name" binding:"required"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type LoginResponse struct {
	Token string `json:"token"`
	User  User   `json:"user"` 
}

type SearchTripsRequest struct {
	FromPlaceID string `form:"fromPlaceId" binding:"required"` 
	ToPlaceID   string `form:"toPlaceId" binding:"required"`
	Date        string `form:"date" binding:"required"` 
}

type CreateBookingRequest struct {
	TripID     string            `json:"tripId" binding:"required"`
	Passengers []PassengerInfo `json:"passengers" binding:"required,dive"` 
}

type MockPaymentRequest struct {
	BookingID string `json:"bookingId" binding:"required"`
}
