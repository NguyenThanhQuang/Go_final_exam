package services

import (
	"context"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/Go_final_exam/bus-booking-backend/src/config"
	"github.com/Go_final_exam/bus-booking-backend/src/models"
	"github.com/Go_final_exam/bus-booking-backend/src/utils"
)

type RegisterInput struct {
	Email    string `json:"email" validate:"required,email"`
	Phone    string `json:"phone" validate:"required"`
	Password string `json:"password" validate:"required,min=6"`
	Name     string `json:"name" validate:"required"`
}

type LoginInput struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

func Register(input RegisterInput) (*models.User, error) {
	userCollection := config.DB.Collection("users")

	count, err := userCollection.CountDocuments(context.TODO(), bson.M{"email": input.Email})
	if err != nil {
		return nil, errors.New("lỗi khi kiểm tra email")
	}
	if count > 0 {
		return nil, errors.New("email đã được sử dụng")
	}

	hashedPassword, err := utils.HashPassword(input.Password)
	if err != nil {
		return nil, errors.New("không thể băm mật khẩu")
	}

	newUser := models.User{
		ID:           primitive.NewObjectID(),
		Email:        input.Email,
		Phone:        input.Phone,
		PasswordHash: hashedPassword,
		Name:         input.Name,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	_, err = userCollection.InsertOne(context.TODO(), newUser)
	if err != nil {
		return nil, errors.New("không thể tạo người dùng")
	}

	return &newUser, nil
}

func Login(input LoginInput) (string, error) {
	userCollection := config.DB.Collection("users")
	var user models.User

	err := userCollection.FindOne(context.TODO(), bson.M{"email": input.Email}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return "", errors.New("email hoặc mật khẩu không chính xác")
		}
		return "", errors.New("lỗi hệ thống khi tìm kiếm người dùng")
	}

	err = utils.CheckPasswordHash(input.Password, user.PasswordHash)
	if err != nil {
		return "", errors.New("email hoặc mật khẩu không chính xác")
	}

	token, err := utils.GenerateToken(user.ID)
	if err != nil {
		return "", errors.New("không thể tạo token xác thực")
	}

	return token, nil
}