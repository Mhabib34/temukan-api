package usecase

import (
	"context"
	"temukan-api/internal/dto"
)

type UserUsecase interface {
	Create(ctx context.Context, request *dto.RegisterRequest) (*dto.UserResponse, error)
}
