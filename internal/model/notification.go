package model

import (
	"time"

	"github.com/google/uuid"
)

type Notification struct {
	ID        uuid.UUID  `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	UserID    uuid.UUID  `gorm:"type:uuid;not null"                             json:"user_id"`
	ReportID  *uuid.UUID `gorm:"type:uuid"                                      json:"report_id"`
	MatchID   *uuid.UUID `gorm:"type:uuid"                                      json:"match_id"`
	Message   string     `gorm:"type:text;not null"                             json:"message"`
	IsRead    bool       `gorm:"not null;default:false"                         json:"is_read"`
	CreatedAt time.Time  `gorm:"not null;default:now()"                         json:"created_at"`

	User  User   `gorm:"foreignKey:UserID"  json:"-"`
	Match *Match `gorm:"foreignKey:MatchID" json:"match,omitempty"`
}
