package repositories

import (
	"bus_booking_api/internal/models"
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type VehicleRepository interface {
	CreateVehicle(ctx context.Context, vehicle *models.Vehicle) (*models.Vehicle, error)
	GetVehicleByID(ctx context.Context, id primitive.ObjectID) (*models.Vehicle, error)
	GetVehiclesByCompanyID(ctx context.Context, companyID primitive.ObjectID) ([]models.Vehicle, error)
}

type mongoVehicleRepository struct {
	collection *mongo.Collection
}

func NewMongoVehicleRepository(db *mongo.Database, collectionName string) VehicleRepository {
	return &mongoVehicleRepository{
		collection: db.Collection(collectionName),
	}
}

func (r *mongoVehicleRepository) CreateVehicle(ctx context.Context, vehicle *models.Vehicle) (*models.Vehicle, error) {
	vehicle.ID = primitive.NewObjectID()
	vehicle.CreatedAt = time.Now()
	vehicle.UpdatedAt = time.Now()

	_, err := r.collection.InsertOne(ctx, vehicle)
	if err != nil {
		return nil, err
	}
	return vehicle, nil
}

func (r *mongoVehicleRepository) GetVehicleByID(ctx context.Context, id primitive.ObjectID) (*models.Vehicle, error) {
	var vehicle models.Vehicle
	err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&vehicle)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return &vehicle, nil
}

func (r *mongoVehicleRepository) GetVehiclesByCompanyID(ctx context.Context, companyID primitive.ObjectID) ([]models.Vehicle, error) {
	var vehicles []models.Vehicle
	cursor, err := r.collection.Find(ctx, bson.M{"companyId": companyID})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	if err = cursor.All(ctx, &vehicles); err != nil {
		return nil, err
	}
	if vehicles == nil { 
		return []models.Vehicle{}, nil 
	}
	return vehicles, nil
}