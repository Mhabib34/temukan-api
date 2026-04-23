package usecase

import (
	"context"
	"temukan-api/internal/dto"
	"temukan-api/internal/exception"
	"temukan-api/internal/helper"
	"temukan-api/internal/model"
	"temukan-api/internal/repository"

	"github.com/go-playground/validator/v10"
)

type UserUsecaseImpl struct {
	repo     repository.UserRepository
	Validate *validator.Validate
}

func NewUserUsecase(repo repository.UserRepository, validate *validator.Validate) UserUsecase {
	return &UserUsecaseImpl{
		repo:     repo,
		Validate: validate,
	}
}

func (u *UserUsecaseImpl) Create(ctx context.Context, request *dto.RegisterRequest) (*dto.UserResponse, error) {
	err := u.Validate.Struct(request)
	exception.PanicIfError(err)

	_, err = u.repo.FIndByEmail(ctx, request.Email)
	if err == nil {
		return &dto.UserResponse{}, exception.NewConflictError("Email already exists")
	}

	hasPassword, err := helper.HashPassword(request.Password)
	exception.PanicIfError(err)

	user := &model.User{
		Name:     request.Name,
		Email:    request.Email,
		Password: hasPassword,
		Role:     request.Role,
	}

	if request.Phone != "" {
		user.Phone = &request.Phone
	}

	result, err := u.repo.Create(ctx, user)
	exception.PanicIfError(err)

	return helper.ToUserResponse(*result), nil
}
