package model

import (
	"time"

	"github.com/google/uuid"
)

type Match struct {
	ID              uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	FoundReportID   uuid.UUID `gorm:"type:uuid;not null"                             json:"found_report_id"`
	MissingReportID uuid.UUID `gorm:"type:uuid;not null"                             json:"missing_report_id"`
	Score           int       `gorm:"not null"                                       json:"score"`
	Notified        bool      `gorm:"not null;default:false"                         json:"notified"`
	CreatedAt       time.Time `gorm:"not null;default:now()"                         json:"created_at"`

	FoundReport   Report `gorm:"foreignKey:FoundReportID"   json:"found_report"`
	MissingReport Report `gorm:"foreignKey:MissingReportID" json:"missing_report"`
}
