package dto

import (
	"time"

	"github.com/google/uuid"
)

type NotificationResponse struct {
	ID        uuid.UUID  `json:"id"`
	Message   string     `json:"message"`
	IsRead    bool       `json:"is_read"`
	ReportID  *uuid.UUID `json:"report_id"`
	MatchID   *uuid.UUID `json:"match_id"`
	CreatedAt time.Time  `json:"created_at"`
}

type GetNotificationsQuery struct {
	IsRead *bool `form:"is_read"`
	Page   int   `form:"page"    binding:"omitempty,min=1"`
	Limit  int   `form:"limit"   binding:"omitempty,min=1,max=100"`
}

type NotificationListData struct {
	Notifications []NotificationResponse `json:"notifications"`
	UnreadCount   int64                  `json:"unread_count"`
	Meta          Pagination             `json:"meta"`
}
