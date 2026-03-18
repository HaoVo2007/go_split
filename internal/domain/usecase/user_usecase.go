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

	userRes "go-split/internal/interface/http/dto/response/user"

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
	GetDashboardSummary(ctx context.Context) (*userRes.DashboardSummaryResponse, error)
}

type userUseCase struct {
	userRepository     repository.UserRepository
	groupRepository    repository.GroupRepository
	expenseRepository  repository.ExpenseRepository
	splitRepository    repository.ExpenseSplitRepository
	cloudinaryUploader *helper.CloudinaryUploader
}

func NewUserUseCase(
	userRepository repository.UserRepository,
	groupRepository repository.GroupRepository,
	expenseRepository repository.ExpenseRepository,
	splitRepository repository.ExpenseSplitRepository,
	cloudinaryUploader *helper.CloudinaryUploader,
) UserUseCase {
	return &userUseCase{
		userRepository:     userRepository,
		groupRepository:    groupRepository,
		expenseRepository:  expenseRepository,
		splitRepository:    splitRepository,
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

	if req.ImageType == "preset" {
		if userCheck.Profile.ImagePublicID != nil {
			u.cloudinaryUploader.DeleteImage(ctx, *userCheck.Profile.ImagePublicID)
		}

		userCheck.Profile.Image = &req.ImageUrl
		userCheck.Profile.ImagePublicID = nil
	}

	if req.ImageType == "upload" && req.Image != nil {
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

func (u *userUseCase) GetDashboardSummary(ctx context.Context) (*userRes.DashboardSummaryResponse, error) {
	userID, err := helper.GetUserID(ctx)
	if err != nil {
		return nil, err
	}

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

	groups, err := u.groupRepository.GetGroups(ctx, userID)
	if err != nil {
		return nil, err
	}

	if len(groups) == 0 {
		return &userRes.DashboardSummaryResponse{
			Balance: userRes.BalanceResponse{
				YouOwed: 0,
				YouPaid: 0,
				Balance: 0,
			},
			Overview: userRes.OverviewResponse{
				TotalGroups:       0,
				TotalTransactions: 0,
				TotalFriends:      0,
			},
			Expenses: userRes.ExpensesResponse{
				TotalPaid:   0,
				TotalShared: 0,
			},
			TopStatistics: userRes.TopStatisticsResponse{
				TopGroup:  userRes.TopGroupResponse{},
				TopFriend: userRes.TopFriendResponse{},
			},
		}, nil
	}

	groupIDs := []string{}
	for _, group := range groups {
		groupIDs = append(groupIDs, group.ID.Hex())
	}

	expenses, err := u.expenseRepository.GetExpensesByGroupIDs(ctx, groupIDs)
	if err != nil {
		return nil, err
	}

	if len(expenses) == 0 {
		return &userRes.DashboardSummaryResponse{
			Balance: userRes.BalanceResponse{
				YouOwed: 0,
				YouPaid: 0,
				Balance: 0,
			},
			Overview: userRes.OverviewResponse{
				TotalGroups:       len(groups),
				TotalTransactions: 0,
				TotalFriends:      0,
			},
			Expenses: userRes.ExpensesResponse{
				TotalPaid:   0,
				TotalShared: 0,
			},
			TopStatistics: userRes.TopStatisticsResponse{
				TopGroup:  userRes.TopGroupResponse{},
				TopFriend: userRes.TopFriendResponse{},
			},
		}, nil
	}

	expenseIDs := []string{}
	for _, expense := range expenses {
		expenseIDs = append(expenseIDs, expense.ID.Hex())
	}

	splits, err := u.splitRepository.GetExpenseSplitsByExpenseIDs(ctx, expenseIDs)
	if err != nil {
		return nil, err
	}

	var totalPaid float64
	var totalShared float64

	for _, expense := range expenses {
		for _, paidByID := range expense.PaidBy {
			if paidByID == userID {
				totalPaid += expense.Amount
			}
		}
	}

	for _, split := range splits {
		if split.UserId == userID {
			totalShared += split.Amount
		}
	}

	youPaid := totalPaid
	youOwed := totalShared
	
	balance := youPaid - youOwed

	overview := userRes.OverviewResponse{
		TotalGroups:       len(groups),
		TotalTransactions: len(expenses),
		TotalFriends:      countUniqueFriends(groups, userID),
	}

	groupMap := make(map[string]float64)
	for _, expense := range expenses {
		groupMap[expense.GroupID] += expense.Amount
	}

	var toGroupID string
	var max float64
	for gid, amount := range groupMap {
		if amount > max {
			max = amount
			toGroupID = gid
		}
	}

	groupID, err := primitive.ObjectIDFromHex(toGroupID)
	if err != nil {
		return nil, err
	}
	groupDetail, err := u.groupRepository.GetGroupById(ctx, groupID)
	if err != nil {
		return nil, err
	}
	if groupDetail == nil {
		return nil, errors.New("group not found")
	}
	
	topGroup := userRes.TopGroupResponse{
		ID: groupDetail.ID.Hex(),
		Name: groupDetail.Name,
		TotalSpend: max,
	}

	friendMap := make(map[string]int)
	for _, expense := range expenses {
		for _, p := range expense.Participants {
			if p != userID {
				friendMap[p]++
			}
		}
	}

	var topFriendID string
	var maxFriend int
	for fid, count := range friendMap {
		if count > maxFriend {
			maxFriend = count
			topFriendID = fid
		}
	}

	friendID, err := primitive.ObjectIDFromHex(topFriendID)
	if err != nil {
		return nil, err
	}

	friendDetail, err := u.userRepository.FindUserByID(ctx, friendID)
	if err != nil {
		return nil, err
	}

	if friendDetail == nil {
		return nil, errors.New("friend not found")
	}
	
	var friendName string
	if friendDetail.Profile != nil && friendDetail.Profile.Name != nil {
		friendName = *friendDetail.Profile.Name
	} else {
		friendName = friendDetail.Email
	}

	topFriend := userRes.TopFriendResponse{
		ID: friendDetail.ID.Hex(),
		Name: friendName,
		TotalTransactions: float64(maxFriend),
	}

	return &userRes.DashboardSummaryResponse{
		Balance: userRes.BalanceResponse{
			YouOwed: youOwed,
			YouPaid: youPaid,
			Balance: balance,
		},
		Overview: overview,
		Expenses: userRes.ExpensesResponse{
			TotalPaid: totalPaid,
			TotalShared: totalShared,
		},
		TopStatistics: userRes.TopStatisticsResponse{
			TopGroup: topGroup,
			TopFriend: topFriend,
		},
	}, nil
}

func countUniqueFriends(groups []*entity.Groups, userID string) int {
	friendIDs := make(map[string]bool)
	for _, group := range groups {
		for _, memberID := range group.Members {
			if memberID != userID {
				friendIDs[memberID] = true
			}
		}
	}
	return len(friendIDs)
}
