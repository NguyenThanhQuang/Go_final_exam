package repositories

import (
	"bus_booking_api/internal/models"
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type TripRepository interface {
	CreateTrip(ctx context.Context, trip *models.Trip) (*models.Trip, error) 
	GetTripByID(ctx context.Context, id primitive.ObjectID) (*models.Trip, error)
	SearchTrips(ctx context.Context, fromLocationName string, toLocationName string, date time.Time) ([]models.Trip, error)
	UpdateTripSeatStatus(ctx context.Context, tripID primitive.ObjectID, seatNumber string, newStatus string, bookingID *primitive.ObjectID) error
}

type mongoTripRepository struct {
	collection *mongo.Collection
}

func NewMongoTripRepository(db *mongo.Database, collectionName string) TripRepository {
	return &mongoTripRepository{
		collection: db.Collection(collectionName),
	}
}

func (r *mongoTripRepository) CreateTrip(ctx context.Context, trip *models.Trip) (*models.Trip, error) {
	trip.ID = primitive.NewObjectID()
	trip.CreatedAt = time.Now()
	trip.UpdatedAt = time.Now()

	for i := range trip.Seats {
		if trip.Seats[i].Status == "" {
			trip.Seats[i].Status = models.SeatStatusAvailable
		}
	}

	_, err := r.collection.InsertOne(ctx, trip)
	if err != nil {
		return nil, err
	}
	return trip, nil
}

func (r *mongoTripRepository) GetTripByID(ctx context.Context, id primitive.ObjectID) (*models.Trip, error) {
	var trip models.Trip
	err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&trip)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return &trip, nil
}

func (r *mongoTripRepository) SearchTrips(ctx context.Context, fromLocationName string, toLocationName string, date time.Time) ([]models.Trip, error) {
	var trips []models.Trip

	startOfDay := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location())
	endOfDay := startOfDay.Add(24 * time.Hour)

	filter := bson.M{
		"route.from.name": bson.M{"$regex": primitive.Regex{Pattern: fromLocationName, Options: "i"}},
		"route.to.name":   bson.M{"$regex": primitive.Regex{Pattern: toLocationName, Options: "i"}},
		"departureTime": bson.M{
			"$gte": startOfDay,
			"$lt":  endOfDay,
		},
		"status": models.TripStatusScheduled, 
	}

	findOptions := options.Find()
	findOptions.SetSort(bson.D{{"departureTime", 1}}) 

	cursor, err := r.collection.Find(ctx, filter, findOptions)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	if err = cursor.All(ctx, &trips); err != nil {
		return nil, err
	}

	if len(trips) == 0 {
		return []models.Trip{}, nil 
	}

	return trips, nil
}

func (r *mongoTripRepository) UpdateTripSeatStatus(ctx context.Context, tripID primitive.ObjectID, seatNumber string, newStatus string, bookingID *primitive.ObjectID) error {
	filter := bson.M{
		"_id":           tripID,
		"seats.seatNumber": seatNumber,
	}
	update := bson.M{
		"$set": bson.M{
			"seats.$.status":    newStatus,
			"seats.$.bookingId": bookingID, 
			"updatedAt":      time.Now(),
		},
	}

	result, err := r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}
	if result.MatchedCount == 0 {
		return mongo.ErrNoDocuments 
	}
	return nil
}