package helper

import (
	"temukan-api/internal/dto"
	"temukan-api/internal/model"
)

func ToUserResponse(user model.User) *dto.UserResponse {
	return &dto.UserResponse{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		Phone:     user.Phone,
		Role:      user.Role,
		CreatedAt: user.CreatedAt,
	}
}
