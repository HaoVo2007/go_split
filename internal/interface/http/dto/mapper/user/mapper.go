package user

import (
	"go-split/internal/domain/entity"
	userRes "go-split/internal/interface/http/dto/response/user"
)

func ToUserResponse(user *entity.Users) *userRes.UserResponse {
	if user.Profile == nil {
		return &userRes.UserResponse{
			ID:      user.ID.Hex(),
			Email:   user.Email,
			Role:    user.Role,
			Name:    nil,
			Image:   nil,
			Address: nil,
			Phone:   nil,
		}
	}

	name := ""
	if user.Profile.Name != nil {
		name = *user.Profile.Name
	}

	image := ""
	if user.Profile.Image != nil {
		image = *user.Profile.Image
	}

	address := ""
	if user.Profile.Address != nil {
		address = *user.Profile.Address
	}

	phone := ""
	if user.Profile.Phone != nil {
		phone = *user.Profile.Phone
	}

	imagePublicID := ""
	if user.Profile.ImagePublicID != nil {
		imagePublicID = *user.Profile.ImagePublicID
	}

	return &userRes.UserResponse{
		ID:            user.ID.Hex(),
		Email:         user.Email,
		Role:          user.Role,
		Name:          &name,
		Image:         &image,
		ImagePublicID: &imagePublicID,
		Address:       &address,
		Phone:         &phone,
	}
}
