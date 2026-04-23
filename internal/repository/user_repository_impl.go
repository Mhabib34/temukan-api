package repository

import (
	"context"
	"temukan-api/internal/exception"
	"temukan-api/internal/model"

	"gorm.io/gorm"
)

type UserRepositoryImpl struct {
	DB *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &UserRepositoryImpl{
		DB: db,
	}
}

func (u *UserRepositoryImpl) Create(ctx context.Context, user *model.User) (*model.User, error) {
	err := u.DB.WithContext(ctx).Create(user).Error
	exception.PanicIfError(err)
	return user, nil
}

func (u *UserRepositoryImpl) FIndByEmail(ctx context.Context, email string) (*model.User, error) {
	var user model.User

	err := u.DB.WithContext(ctx).
		Where("email = ?", email).
		First(&user).Error

	if err != nil {
		return nil, err
	}
	return &user, nil
}
