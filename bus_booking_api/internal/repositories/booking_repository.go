package repositories

import (
	"bus_booking_api/internal/models"
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type BookingRepository interface {
	CreateBooking(ctx context.Context, booking *models.Booking) (*models.Booking, error)
	GetBookingByID(ctx context.Context, id primitive.ObjectID) (*models.Booking, error)
	GetBookingByTicketCode(ctx context.Context, ticketCode string) (*models.Booking, error)
	UpdateBookingStatus(ctx context.Context, id primitive.ObjectID, status string, paymentStatus *string, paymentGatewayTxID *string, ticketCode *string) error
}

type mongoBookingRepository struct {
	collection *mongo.Collection
}

func NewMongoBookingRepository(db *mongo.Database, collectionName string) BookingRepository {
	return &mongoBookingRepository{
		collection: db.Collection(collectionName),
	}
}

func (r *mongoBookingRepository) CreateBooking(ctx context.Context, booking *models.Booking) (*models.Booking, error) {
	booking.ID = primitive.NewObjectID()
	booking.BookingTime = time.Now()
	booking.CreatedAt = time.Now()
	booking.UpdatedAt = time.Now()

	_, err := r.collection.InsertOne(ctx, booking)
	if err != nil {
		return nil, err
	}
	return booking, nil
}

func (r *mongoBookingRepository) GetBookingByID(ctx context.Context, id primitive.ObjectID) (*models.Booking, error) {
	var booking models.Booking
	err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&booking)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return &booking, nil
}

func (r *mongoBookingRepository) GetBookingByTicketCode(ctx context.Context, ticketCode string) (*models.Booking, error) {
	var booking models.Booking
	err := r.collection.FindOne(ctx, bson.M{"ticketCode": ticketCode}).Decode(&booking)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return &booking, nil
}

func (r *mongoBookingRepository) UpdateBookingStatus(ctx context.Context, id primitive.ObjectID, status string, paymentStatus *string, paymentGatewayTxID *string, ticketCode *string) error {
	updateFields := bson.M{
		"status":    status,
		"updatedAt": time.Now(),
	}
	if paymentStatus != nil {
		updateFields["paymentStatus"] = *paymentStatus
	}
	if paymentGatewayTxID != nil {
		updateFields["paymentGatewayTransactionId"] = *paymentGatewayTxID
	}
	if ticketCode != nil {
		updateFields["ticketCode"] = *ticketCode
	}
	if status == models.BookingStatusHeld {
		var booking models.Booking
		err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&booking)
		if err == nil && booking.HeldUntil == nil {
			heldUntil := time.Now().Add(15 * time.Minute)
			updateFields["heldUntil"] = heldUntil
		}
	}

	update := bson.M{"$set": updateFields}

	result, err := r.collection.UpdateOne(ctx, bson.M{"_id": id}, update)
	if err != nil {
		return err
	}
	if result.MatchedCount == 0 {
		return mongo.ErrNoDocuments 
	}
	return nil
}