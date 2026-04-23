package dto

import (
	"temukan-api/internal/model"
	"time"

	"github.com/google/uuid"
)

type RegisterRequest struct {
	Name     string     `json:"name"     binding:"required"`
	Email    string     `json:"email"    binding:"required,email"`
	Password string     `json:"password" binding:"required"`
	Role     model.Role `json:"role"     binding:"required,oneof=finder seeker volunteer"`
	Phone    string     `json:"phone"`
}

type UserResponse struct {
	ID        uuid.UUID  `json:"id"`
	Name      string     `json:"name"`
	Email     string     `json:"email"`
	Role      model.Role `json:"role"`
	Phone     *string    `json:"phone"`
	CreatedAt time.Time  `json:"created_at"`
}
