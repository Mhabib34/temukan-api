package model

import (
	"time"

	"github.com/google/uuid"
)

type Role string

const (
	RoleFinder    Role = "finder"
	RoleSeeker    Role = "seeker"
	RoleVolunteer Role = "volunteer"
)

// User adalah entitas utama untuk autentikasi
type User struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	Name      string    `gorm:"type:varchar(255);not null" json:"name"`
	Email     string    `gorm:"type:varchar(255);uniqueIndex;not null" json:"email"`
	Password  string    `gorm:"type:varchar(255);not null" json:"-"`
	Role      Role      `gorm:"type:varchar(50);not null" json:"role"`
	Phone     *string   `gorm:"type:varchar(20)" json:"phone"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// RefreshToken menyimpan refresh token yang aktif di DB
type RefreshToken struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	UserID    uuid.UUID `gorm:"type:uuid;not null;index"`
	Token     string    `gorm:"type:text;uniqueIndex;not null"`
	ExpiresAt time.Time `gorm:"not null"`
	CreatedAt time.Time
}

func (RefreshToken) TableName() string {
	return "refresh_tokens"
}
