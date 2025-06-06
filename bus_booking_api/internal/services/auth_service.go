package services

import (
	"bus_booking_api/internal/config"
	"bus_booking_api/internal/models"
	"bus_booking_api/internal/repositories"
	"bus_booking_api/internal/utils"
	"context"
	"errors"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type AuthService interface {
	Register(ctx context.Context, name, email, phone, password string) (*models.User, error)
	Login(ctx context.Context, email, password string) (user *models.User, token string, err error)
}

type authService struct {
	userRepo repositories.UserRepository
	cfg      *config.Config
}

func NewAuthService(userRepo repositories.UserRepository, cfg *config.Config) AuthService {
	return &authService{
		userRepo: userRepo,
		cfg:      cfg,
	}
}

func (s *authService) Register(ctx context.Context, name, email, phone, password string) (*models.User, error) {
	if name == "" || email == "" || phone == "" || password == "" {
		return nil, errors.New("tên, email, số điện thoại và mật khẩu không được để trống")
	}
	if len(password) < 6 {
		return nil, errors.New("mật khẩu phải có ít nhất 6 ký tự")
	}
	if !strings.Contains(email, "@") {
		return nil, errors.New("địa chỉ email không hợp lệ")
	}

	existingUserByEmail, err := s.userRepo.GetUserByEmail(ctx, email)
	if err != nil && err != mongo.ErrNoDocuments { 
		return nil, errors.New("lỗi khi kiểm tra email: " + err.Error())
	}
	if existingUserByEmail != nil {
		return nil, errors.New("email đã được sử dụng")
	}

	hashedPassword, err := utils.HashPassword(password)
	if err != nil {
		return nil, errors.New("lỗi khi mã hóa mật khẩu: " + err.Error())
	}

	newUser := &models.User{
		Name:         name,
		Email:        strings.ToLower(email),
		Phone:        phone,
		PasswordHash: hashedPassword,
		Role:         models.RoleUser, 
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	createdUser, err := s.userRepo.CreateUser(ctx, newUser)
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			if strings.Contains(err.Error(), "email_1") { 
				return nil, errors.New("email đã tồn tại")
			}
			if strings.Contains(err.Error(), "phone_1") { 
				return nil, errors.New("số điện thoại đã tồn tại")
			}
			return nil, errors.New("người dùng với thông tin này đã tồn tại")
		}
		return nil, errors.New("không thể tạo người dùng: " + err.Error())
	}

	createdUser.PasswordHash = ""
	return createdUser, nil
}

func (s *authService) Login(ctx context.Context, email, password string) (user *models.User, token string, err error) {
	if email == "" || password == "" {
		return nil, "", errors.New("email và mật khẩu không được để trống")
	}

	foundUser, err := s.userRepo.GetUserByEmail(ctx, strings.ToLower(email))
	if err != nil {
		return nil, "", errors.New("lỗi khi tìm kiếm người dùng: " + err.Error())
	}
	if foundUser == nil {
		return nil, "", errors.New("email hoặc mật khẩu không chính xác")
	}

	if !utils.CheckPasswordHash(password, foundUser.PasswordHash) {
		return nil, "", errors.New("email hoặc mật khẩu không chính xác")
	}

	var companyIDStr string
	if foundUser.CompanyID != nil && *foundUser.CompanyID != primitive.NilObjectID {
		companyIDStr = foundUser.CompanyID.Hex()
	}

	jwtToken, err := utils.GenerateToken(foundUser.ID.Hex(), foundUser.Email, foundUser.Role, companyIDStr, s.cfg)
	if err != nil {
		return nil, "", errors.New("lỗi khi tạo token xác thực: " + err.Error())
	}

	foundUser.PasswordHash = ""
	return foundUser, jwtToken, nil
}