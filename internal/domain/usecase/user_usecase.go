package usecase

import (
	"context"
	"errors"
	"fmt"
	"go-split/internal/domain/entity"
	"go-split/internal/domain/repository"
	"go-split/internal/interface/http/dto/request/user"
	"go-split/pkg/libs/contants"
	"go-split/pkg/libs/helper"
	"log"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

type UserUseCase interface {
	CreateUser(ctx context.Context, req user.CreateUserRequest) error
	LoginUser(ctx context.Context, req user.LoginUserRequest) (*entity.Users, error)
	RegisterUser(ctx context.Context, req user.RegisterUserRequest) (*entity.Users, error)
	LogoutUser(ctx context.Context) error
	UpdateProfile(ctx context.Context, req user.UpdateProfileRequest) error
	GetUsers(ctx context.Context) ([]*entity.Users, error)
	GetUserById(ctx context.Context, userID string) (*entity.Users, error)
	UpdateUser(ctx context.Context, userID string, req user.UpdateUserRequest) error
	SoftDeleteUser(ctx context.Context, userID string) error
	RefreshToken(ctx context.Context, refreshToken string) (string, string, error)
}

type userUseCase struct {
	userRepository     repository.UserRepository
	cloudinaryUploader *helper.CloudinaryUploader
}

func NewUserUseCase(
	userRepository repository.UserRepository,
	cloudinaryUploader *helper.CloudinaryUploader,
) UserUseCase {
	return &userUseCase{
		userRepository:     userRepository,
		cloudinaryUploader: cloudinaryUploader,
	}
}

func (u *userUseCase) CreateUser(ctx context.Context, req user.CreateUserRequest) error {
	userCheck, err := u.userRepository.FindUserByEmail(ctx, req.Email)
	if err != nil {
		return err
	}
	if userCheck != nil {
		return errors.New("user already exists")
	}

	hashedPassword := u.HashPassword(req.Password)
	user := entity.Users{
		ID:           primitive.NewObjectID(),
		Email:        req.Email,
		Password:     hashedPassword,
		Role:         req.Role,
		Status:       contants.StatusActive,
		Token:        "",
		RefreshToken: "",
		IsDeleted:    false,
		DeletedAt:    nil,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}
	if err := u.userRepository.CreateUser(ctx, user); err != nil {
		return err
	}

	return nil
}

func (u *userUseCase) LoginUser(ctx context.Context, req user.LoginUserRequest) (*entity.Users, error) {
	userCheck, err := u.userRepository.FindUserByEmail(ctx, req.Email)
	if err != nil {
		return nil, err
	}

	if userCheck == nil {
		return nil, errors.New("invalid email or password")
	}

	isValid, msg := u.VerifyPassword(userCheck.Password, req.Password)
	if !isValid {
		return nil, errors.New(msg)
	}

	token, refreshToken := u.GenerateToken(userCheck)

	updateFields := bson.M{
		"token":         token,
		"refresh_token": refreshToken,
		"updated_at":    time.Now(),
	}

	err = u.userRepository.UpdateUser(ctx, userCheck.ID, updateFields)
	if err != nil {
		return nil, err
	}

	userCheck.Token = token
	userCheck.RefreshToken = refreshToken
	userCheck.Password = ""

	return userCheck, nil
}

func (u *userUseCase) RegisterUser(ctx context.Context, req user.RegisterUserRequest) (*entity.Users, error) {
	userCheck, err := u.userRepository.FindUserByEmail(ctx, req.Email)
	if err != nil {
		return nil, err
	}

	if userCheck != nil {
		return nil, errors.New("user already exists")
	}

	hashedPassword := u.HashPassword(req.Password)
	user := entity.Users{
		ID:           primitive.NewObjectID(),
		Email:        req.Email,
		Password:     hashedPassword,
		Role:         "user",
		Status:       contants.StatusActive,
		Token:        "",
		RefreshToken: "",
		IsDeleted:    false,
		DeletedAt:    nil,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	token, refreshToken := u.GenerateToken(&user)
	user.Token = token
	user.RefreshToken = refreshToken

	if err := u.userRepository.CreateUser(ctx, user); err != nil {
		return nil, err
	}
	user.Password = ""

	return &user, nil
}

func (u *userUseCase) HashPassword(password string) string {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		log.Printf("Failed to hash password: %v", err)
	}
	return string(bytes)
}

func (u *userUseCase) VerifyPassword(userPassword string, providedPassword string) (bool, string) {
	err := bcrypt.CompareHashAndPassword([]byte(userPassword), []byte(providedPassword))
	check := true
	msg := ""

	if err != nil {
		msg = "Email or password is incorrect"
		check = false
	}
	return check, msg
}

func (u *userUseCase) GenerateToken(user *entity.Users) (string, string) {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		log.Print("JWT_SECRET not set")
	}

	claims := jwt.MapClaims{
		"user_id": user.ID.Hex(),
		"role":    user.Role,
		"type":    "access",
		"exp":     jwt.NewNumericDate(time.Now().Add(time.Hour * 800000)),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		log.Printf("Failed to generate access token: %v", err)
	}

	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID.Hex(),
		"role":    user.Role,
		"type":    "refresh",
		"exp":     jwt.NewNumericDate(time.Now().Add(time.Hour * 480000)),
	})
	refreshTokenString, err := refreshToken.SignedString([]byte(secret))
	if err != nil {
		log.Printf("Failed to generate refresh token: %v", err)
	}

	return tokenString, refreshTokenString
}

func (u *userUseCase) LogoutUser(ctx context.Context) error {
	userID, err := helper.GetUserID(ctx)
	if err != nil {
		return err
	}

	objectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return err
	}

	updateFields := bson.M{
		"token":         "",
		"refresh_token": "",
		"updated_at":    time.Now(),
	}

	err = u.userRepository.UpdateUser(ctx, objectID, updateFields)
	if err != nil {
		return err
	}

	return nil
}

func (u *userUseCase) UpdateProfile(ctx context.Context, req user.UpdateProfileRequest) error {
	userID, err := helper.GetUserID(ctx)
	if err != nil {
		return err
	}

	objectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return err
	}

	userCheck, err := u.userRepository.FindUserByID(ctx, objectID)
	if err != nil {
		return err
	}

	if userCheck == nil {
		return errors.New("user not found")
	}

	if userCheck.Profile == nil {
		userCheck.Profile = &entity.Profile{}
	}

	if req.Name != "" {
		userCheck.Profile.Name = &req.Name
	}

	if req.Image != nil {
		if userCheck.Profile.ImagePublicID != nil {
			u.cloudinaryUploader.DeleteImage(ctx, *userCheck.Profile.ImagePublicID)
		}
		tempPath := fmt.Sprintf("temp_%s", req.Image.Filename)
		if err := helper.SaveUploadedFile(req.Image, tempPath); err != nil {
			return err
		}
		defer os.Remove(tempPath)
		image, publicID, err := u.cloudinaryUploader.UploadImage(ctx, tempPath, "users")
		if err != nil {
			return err
		}
		userCheck.Profile.Image = &image
		userCheck.Profile.ImagePublicID = &publicID
	}

	if req.Address != "" {
		userCheck.Profile.Address = &req.Address
	}

	if req.Phone != "" {
		userCheck.Profile.Phone = &req.Phone
	}

	updateFields := bson.M{
		"profile":    userCheck.Profile,
		"updated_at": time.Now(),
	}

	err = u.userRepository.UpdateUser(ctx, userCheck.ID, updateFields)
	if err != nil {
		return err
	}

	return nil
}

func (u *userUseCase) GetUsers(ctx context.Context) ([]*entity.Users, error) {
	users, err := u.userRepository.GetUsers(ctx)
	if err != nil {
		return nil, err
	}
	return users, nil
}

func (u *userUseCase) GetUserById(ctx context.Context, userID string) (*entity.Users, error) {
	objectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, err
	}
	user, err := u.userRepository.FindUserByID(ctx, objectID)
	if err != nil {
		return nil, err
	}

	if user == nil {
		return nil, errors.New("user not found")
	}

	return user, nil
}

func (u *userUseCase) UpdateUser(ctx context.Context, userID string, req user.UpdateUserRequest) error {
	objectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return err
	}

	userCheck, err := u.userRepository.FindUserByID(ctx, objectID)
	if err != nil {
		return err
	}

	if userCheck == nil {
		return errors.New("user not found")
	}

	if req.Email != "" {
		userCheck.Email = req.Email
	}
	if req.Password != "" {
		userCheck.Password = u.HashPassword(req.Password)
	}
	if req.Role != "" {
		userCheck.Role = req.Role
	}

	updateFields := bson.M{
		"email":      userCheck.Email,
		"password":   userCheck.Password,
		"role":       userCheck.Role,
		"updated_at": time.Now(),
	}

	err = u.userRepository.UpdateUser(ctx, userCheck.ID, updateFields)
	if err != nil {
		return err
	}

	return nil
}

func (u *userUseCase) SoftDeleteUser(ctx context.Context, userID string) error {
	objectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return err
	}

	userCheck, err := u.userRepository.FindUserByID(ctx, objectID)
	if err != nil {
		return err
	}

	if userCheck == nil {
		return errors.New("user not found")
	}

	updateFields := bson.M{
		"is_deleted": true,
		"deleted_at": time.Now(),
	}

	err = u.userRepository.UpdateUser(ctx, userCheck.ID, updateFields)
	if err != nil {
		return err
	}
	return nil
}

func (u *userUseCase) RefreshToken(ctx context.Context, refreshToken string) (string, string, error) {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		log.Print("JWT_SECRET not set")
	}

	token, err := jwt.Parse(refreshToken, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})
	if err != nil || !token.Valid {
		return "", "", errors.New("invalid refresh token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", "", errors.New("invalid token claims")
	}

	if claims["type"] != "refresh" {
		return "", "", errors.New("invalid token type")
	}

	userID, ok := claims["user_id"].(string)
	if !ok {
		return "", "", errors.New("invalid user ID in token")
	}

	objectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return "", "", errors.New("invalid user ID format")
	}

	user, err := u.userRepository.FindUserByID(ctx, objectID)
	if err != nil {
		return "", "", errors.New("failed to find user: " + err.Error())
	}

	if user == nil {
		return "", "", errors.New("user not found")
	}

	if user.RefreshToken != refreshToken {
		return "", "", errors.New("invalid refresh token")
	}

	accessToken, refreshToken := u.GenerateToken(user)

	updateFields := bson.M{
		"token":         accessToken,
		"refresh_token": refreshToken,
		"updated_at":    time.Now(),
	}

	err = u.userRepository.UpdateUser(ctx, user.ID, updateFields)
	if err != nil {
		return "", "", errors.New("failed to update user: " + err.Error())
	}

	return accessToken, refreshToken, nil
}
