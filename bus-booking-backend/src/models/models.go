package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID           primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Email        string             `json:"email" bson:"email" validate:"required,email"`
	Phone        string             `json:"phone" bson:"phone" validate:"required"`
	PasswordHash string             `json:"-" bson:"passwordHash"` 
	Name         string             `json:"name" bson:"name" validate:"required"`
	CreatedAt    time.Time          `json:"createdAt" bson:"createdAt"`
	UpdatedAt    time.Time          `json:"updatedAt" bson:"updatedAt"`
}

type Company struct {
	ID          primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Name        string             `json:"name" bson:"name"`
	Code        string             `json:"code" bson:"code"`
	Description string             `json:"description,omitempty" bson:"description,omitempty"`
	LogoURL     string             `json:"logoUrl,omitempty" bson:"logoUrl,omitempty"`
}

type Vehicle struct {
	ID         primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Type       string             `json:"type" bson:"type"`
	TotalSeats int                `json:"totalSeats" bson:"totalSeats"`
	SeatMap    interface{}        `json:"seatMap,omitempty" bson:"seatMap,omitempty"` 
}

type Trip struct {
	ID                  primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	CompanyID           primitive.ObjectID `json:"companyId" bson:"companyId"`
	VehicleID           primitive.ObjectID `json:"vehicleId" bson:"vehicleId"`
	Route               Route              `json:"route" bson:"route"`
	DepartureTime       time.Time          `json:"departureTime" bson:"departureTime"`
	ExpectedArrivalTime time.Time          `json:"expectedArrivalTime" bson:"expectedArrivalTime"`
	Price               float64            `json:"price" bson:"price"`
	Seats               []Seat             `json:"seats" bson:"seats"`
	AvailableSeats      int                `json:"availableSeats" bson:"-"`
	CompanyInfo *Company `json:"companyInfo,omitempty" bson:"-"` 
	VehicleInfo *Vehicle `json:"vehicleInfo,omitempty" bson:"-"` 
}

type Seat struct {
	SeatNumber string `json:"seatNumber" bson:"seatNumber"`
	Status     string `json:"status" bson:"status"` 
}

type Route struct {
	From LocationPoint `json:"from" bson:"from"`
	To   LocationPoint `json:"to" bson:"to"`
}

type LocationPoint struct {
	Name string `json:"name" bson:"name"`
}

type Booking struct {
	ID          primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	UserID      primitive.ObjectID `json:"userId" bson:"userId"`
	TripID      primitive.ObjectID `json:"tripId" bson:"tripId"`
	BookingTime time.Time          `json:"bookingTime" bson:"bookingTime"`
	Status      string             `json:"status" bson:"status"` 
	TotalAmount float64            `json:"totalAmount" bson:"totalAmount"`
	Passengers  []Passenger        `json:"passengers" bson:"passengers"`
	TicketCode  string             `json:"ticketCode,omitempty" bson:"ticketCode,omitempty"`
	CreatedAt   time.Time          `json:"createdAt" bson:"createdAt"`
	UpdatedAt   time.Time          `json:"updatedAt" bson:"updatedAt"`

	TripInfo *Trip `json:"tripInfo,omitempty" bson:"-"`
}

type Passenger struct {
	Name       string `json:"name" bson:"name"`
	Phone      string `json:"phone" bson:"phone"`
	SeatNumber string `json:"seatNumber" bson:"seatNumber"`
}