package usecase

import (
	"context"
	"errors"
	"fmt"
	"go-split/internal/domain/entity"
	"go-split/internal/domain/repository"
	groupMapper "go-split/internal/interface/http/dto/mapper/group"
	"go-split/internal/interface/http/dto/request/group"
	groupRes "go-split/internal/interface/http/dto/response/group"
	"go-split/pkg/libs/helper"
	"os"
	"slices"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type GroupUseCase interface {
	CreateGroup(ctx context.Context, req group.CreateGroupRequest) error
	GetGroups(ctx context.Context) (*groupRes.ListGroupResponse, error)
	GetGroupById(ctx context.Context, id string) (*entity.Groups, error)
	UpdateGroup(ctx context.Context, id string, req group.UpdateGroupRequest) error
	DeleteGroup(ctx context.Context, id string) error

	AddGroupMember(ctx context.Context, id string, req group.AddGroupMemberRequest) error
	RemoveGroupMember(ctx context.Context, id string, memberId string) error
}

type groupUseCase struct {
	groupRepository        repository.GroupRepository
	userRepository         repository.UserRepository
	expenseRepository      repository.ExpenseRepository
	expenseSplitRepository repository.ExpenseSplitRepository
	cloudinaryUploader     *helper.CloudinaryUploader
}

func NewGroupUseCase(
	groupRepository repository.GroupRepository,
	userRepository repository.UserRepository,
	expenseRepository repository.ExpenseRepository,
	expenseSplitRepository repository.ExpenseSplitRepository,
	cloudinaryUploader *helper.CloudinaryUploader,
) GroupUseCase {
	return &groupUseCase{
		groupRepository:    groupRepository,
		userRepository:     userRepository,
		expenseRepository:  expenseRepository,
		expenseSplitRepository: expenseSplitRepository,
		cloudinaryUploader: cloudinaryUploader,
	}
}

func (u *groupUseCase) CreateGroup(ctx context.Context, req group.CreateGroupRequest) error {
	var imageURL string
	var imagePublicID string
	if req.Image != nil {
		tempPath := fmt.Sprintf("temp_%s", req.Image.Filename)
		if err := helper.SaveUploadedFile(req.Image, tempPath); err != nil {
			return err
		}
		defer os.Remove(tempPath)
		image, publicID, err := u.cloudinaryUploader.UploadImage(ctx, tempPath, "groups")
		if err != nil {
			return err
		}
		imageURL = image
		imagePublicID = publicID
	}

	userID, err := helper.GetUserID(ctx)
	if err != nil {
		return err
	}

	group := entity.Groups{
		ID:            primitive.NewObjectID(),
		Name:          req.Name,
		Description:   req.Description,
		Image:         imageURL,
		ImagePublicID: imagePublicID,
		Members:       []string{},
		IsDeleted:     false,
		CreatedBy:     userID,
		DeletedAt:     nil,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	return u.groupRepository.CreateGroup(ctx, group)
}

func (u *groupUseCase) GetGroups(ctx context.Context) (*groupRes.ListGroupResponse, error) {
	userID, err := helper.GetUserID(ctx)
	if err != nil {
		return nil, err
	}

	groups, err := u.groupRepository.GetGroups(ctx, userID)
	if err != nil {
		return nil, err
	}

	memberIDs := []primitive.ObjectID{}
	for _, group := range groups {
		for _, memberID := range group.Members {
			memberIDObject, err := primitive.ObjectIDFromHex(memberID)
			if err != nil {
				return nil, err
			}
			memberIDs = append(memberIDs, memberIDObject)
		}
	}

	membersUsers, err := u.userRepository.GetUsersByGroupIDs(ctx, memberIDs)
	if err != nil {
		return nil, err
	}

	userMap := make(map[string]*entity.Users)
	for _, user := range membersUsers {
		userMap[user.ID.Hex()] = user
	}

	responses := make([]*groupRes.GroupResponse, len(groups))
	for i, group := range groups {
		var groupMembers []*entity.Users
		for _, memberID := range group.Members {
			if user, ok := userMap[memberID]; ok {
				groupMembers = append(groupMembers, user)
			}
		}
		responses[i] = groupMapper.ToGroupResponse(group, groupMembers)
	}

	totalPaid := 0.0
	totalOwed := 0.0

	for _, group := range groups {
		expenses, err := u.expenseRepository.GetExpensesByGroupID(ctx, group.ID.Hex())
		if err != nil {
			return nil, err
		}
		for _, expense := range expenses {
			share := expense.Amount / float64(len(expense.PaidBy))
			for _, paidByID := range expense.PaidBy {
				if paidByID == userID {
					totalPaid += share
					break
				}
			}
			splits, err := u.expenseSplitRepository.GetExpenseSplitsByExpenseID(ctx, expense.ID.Hex())
			if err != nil {
				return nil, err
			}
			for _, split := range splits {
				if split.UserId == userID {
					totalOwed += split.Amount
				}
			}
		}
	}

	listResponse := &groupRes.ListGroupResponse{
		TotalGroups: len(groups),
		TotalPaid:   totalPaid,
		TotalOwed:   totalOwed,
		Groups:      responses,
	}

	return listResponse, nil
}

func (u *groupUseCase) GetGroupById(ctx context.Context, id string) (*entity.Groups, error) {
	groupID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	group, err := u.groupRepository.GetGroupById(ctx, groupID)
	if err != nil {
		return nil, err
	}
	return group, nil
}

func (u *groupUseCase) UpdateGroup(ctx context.Context, id string, req group.UpdateGroupRequest) error {
	groupID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	group, err := u.groupRepository.GetGroupById(ctx, groupID)
	if err != nil {
		return err
	}

	imageUrl := group.Image
	imagePublicID := group.ImagePublicID
	if req.Image != nil {
		if group.ImagePublicID != "" {
			u.cloudinaryUploader.DeleteImage(ctx, group.ImagePublicID)
		}
		tempPath := fmt.Sprintf("temp_%s", req.Image.Filename)
		if err := helper.SaveUploadedFile(req.Image, tempPath); err != nil {
			return err
		}
		defer os.Remove(tempPath)
		image, publicID, err := u.cloudinaryUploader.UploadImage(ctx, tempPath, "groups")
		if err != nil {
			return err
		}
		imageUrl = image
		imagePublicID = publicID
	}

	group.Name = req.Name
	group.Description = req.Description
	group.Image = imageUrl
	group.ImagePublicID = imagePublicID
	group.UpdatedAt = time.Now()

	return u.groupRepository.UpdateGroup(ctx, groupID, group)
}

func (u *groupUseCase) DeleteGroup(ctx context.Context, id string) error {
	groupID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	group, err := u.groupRepository.GetGroupById(ctx, groupID)
	if err != nil {
		return err
	}

	if group.ImagePublicID != "" {
		u.cloudinaryUploader.DeleteImage(ctx, group.ImagePublicID)
	}

	return u.groupRepository.DeleteGroup(ctx, groupID)
}

func (u *groupUseCase) AddGroupMember(ctx context.Context, id string, req group.AddGroupMemberRequest) error {
	groupID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	user, err := u.userRepository.FindUserByEmail(ctx, req.Email)
	if err != nil {
		return err
	}
	if user == nil {
		return errors.New("user not found")
	}

	group, err := u.groupRepository.GetGroupById(ctx, groupID)
	if err != nil {
		return err
	}

	if slices.Contains(group.Members, user.ID.Hex()) || group.CreatedBy == user.ID.Hex() {
		return errors.New("user already in group")
	}

	group.Members = append(group.Members, user.ID.Hex())
	group.UpdatedAt = time.Now()

	return u.groupRepository.UpdateGroup(ctx, groupID, group)
}

func (u *groupUseCase) RemoveGroupMember(ctx context.Context, id string, memberId string) error {
	groupID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	group, err := u.groupRepository.GetGroupById(ctx, groupID)
	if err != nil {
		return err
	}

	memberID, err := primitive.ObjectIDFromHex(memberId)
	if err != nil {
		return err
	}

	if slices.Contains(group.Members, memberID.Hex()) {
		return errors.New("member not in group")
	}

	group.Members = slices.Delete(group.Members, slices.Index(group.Members, memberID.Hex()), 1)
	group.UpdatedAt = time.Now()

	return u.groupRepository.UpdateGroup(ctx, groupID, group)
}
