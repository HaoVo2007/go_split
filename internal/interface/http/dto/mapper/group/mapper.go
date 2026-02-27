package group

import (
	"go-split/internal/domain/entity"
	groupRes "go-split/internal/interface/http/dto/response/group"
	userRes "go-split/internal/interface/http/dto/response/user"
)

func ToGroupResponse(group *entity.Groups, members []*entity.Users) *groupRes.GroupResponse {
	membersResponse := make([]*userRes.UserResponse, len(members))
	for i, member := range members {
		if member.Profile == nil {
			membersResponse[i] = &userRes.UserResponse{
				ID:      member.ID.Hex(),
				Email:   member.Email,
				Role:    member.Role,
				Name:    nil,
				Image:   nil,
				Address: nil,
				Phone:   nil,
			}
			continue
		}

		name := ""
		if member.Profile.Name != nil {
			name = *member.Profile.Name
		}

		image := ""
		if member.Profile.Image != nil {
			image = *member.Profile.Image
		}

		address := ""
		if member.Profile.Address != nil {
			address = *member.Profile.Address
		}

		phone := ""
		if member.Profile.Phone != nil {
			phone = *member.Profile.Phone
		}
		
		membersResponse[i] = &userRes.UserResponse{
			ID:      member.ID.Hex(),
			Email:   member.Email,
			Role:    member.Role,
			Name:    &name,
			Image:   &image,
			Address: &address,
			Phone:   &phone,
		}
	}

	return &groupRes.GroupResponse{
		ID:            group.ID.Hex(),
		Name:          group.Name,
		Image:         group.Image,
		ImagePublicID: group.ImagePublicID,
		Members:       membersResponse,
		Description:   group.Description,
		TotalMembers:  len(membersResponse),
		CreatedBy:     group.CreatedBy,
		CreatedAt:     group.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:     group.UpdatedAt.Format("2006-01-02 15:04:05"),
	}
}
