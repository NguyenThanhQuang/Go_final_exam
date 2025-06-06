package repositories

import (
	"bus_booking_api/internal/models"
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type CompanyRepository interface {
	CreateCompany(ctx context.Context, company *models.Company) (*models.Company, error)
	GetCompanyByID(ctx context.Context, id primitive.ObjectID) (*models.Company, error)
	GetCompanyByCode(ctx context.Context, code string) (*models.Company, error) 
}

type mongoCompanyRepository struct {
	collection *mongo.Collection
}

func NewMongoCompanyRepository(db *mongo.Database, collectionName string) CompanyRepository {
	return &mongoCompanyRepository{
		collection: db.Collection(collectionName),
	}
}

func (r *mongoCompanyRepository) CreateCompany(ctx context.Context, company *models.Company) (*models.Company, error) {
	company.ID = primitive.NewObjectID()
	company.CreatedAt = time.Now()
	company.UpdatedAt = time.Now()
	if company.IsActive == false && company.ID == primitive.NilObjectID { 
		company.IsActive = true
	}


	_, err := r.collection.InsertOne(ctx, company)
	if err != nil {
		return nil, err
	}
	return company, nil
}

func (r *mongoCompanyRepository) GetCompanyByID(ctx context.Context, id primitive.ObjectID) (*models.Company, error) {
	var company models.Company
	err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&company)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return &company, nil
}

func (r *mongoCompanyRepository) GetCompanyByCode(ctx context.Context, code string) (*models.Company, error) {
	var company models.Company
	err := r.collection.FindOne(ctx, bson.M{"code": code}).Decode(&company)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return &company, nil
}