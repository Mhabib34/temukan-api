package repository

import (
	"context"
	"temukan-api/internal/model"
)

type UserRepository interface {
	Create(ctx context.Context, user *model.User) (*model.User, error)
	FIndByEmail(ctx context.Context, email string) (*model.User, error)
}
