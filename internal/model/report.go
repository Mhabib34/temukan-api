package model

import (
	"time"

	"github.com/google/uuid"
)

type ReportType string
type ReportGender string
type ReportStatus string

const (
	ReportTypeFound   ReportType = "found"
	ReportTypeMissing ReportType = "missing"

	ReportGenderMale    ReportGender = "male"
	ReportGenderFemale  ReportGender = "female"
	ReportGenderUnknown ReportGender = "unknown"

	ReportStatusActive   ReportStatus = "active"
	ReportStatusResolved ReportStatus = "resolved"
)

type Report struct {
	ID               uuid.UUID    `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	ReporterID       uuid.UUID    `gorm:"type:uuid;not null;index"                       json:"reporter_id"`
	Reporter         User         `gorm:"foreignKey:ReporterID;references:ID"            json:"reporter"`
	Type             ReportType   `gorm:"type:report_type;not null"                      json:"type"`
	Name             *string      `gorm:"type:varchar(255)"                              json:"name"`
	Gender           ReportGender `gorm:"type:report_gender;not null;default:unknown"    json:"gender"`
	EstimatedAge     *int         `gorm:"check:estimated_age >= 0 AND estimated_age <= 120" json:"estimated_age"`
	PhotoURL         *string      `gorm:"type:varchar(500)"                              json:"photo_url"`
	Description      string       `gorm:"type:text;not null"                             json:"description"`
	LastSeenLocation string       `gorm:"type:varchar(500);not null"                     json:"last_seen_location"`
	City             string       `gorm:"type:varchar(100);not null"                     json:"city"`
	Province         string       `gorm:"type:varchar(100);not null"                     json:"province"`
	Latitude         *float64     `gorm:"type:double precision"                          json:"latitude"`
	Longitude        *float64     `gorm:"type:double precision"                          json:"longitude"`
	Status           ReportStatus `gorm:"type:report_status;not null;default:active"     json:"status"`
	CreatedAt        time.Time    `                                                      json:"created_at"`
	UpdatedAt        time.Time    `                                                      json:"updated_at"`
}
