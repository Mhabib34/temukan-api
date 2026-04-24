package dto

import (
	"time"

	"github.com/google/uuid"
)

type MatchResponse struct {
	ID            uuid.UUID      `json:"id"`
	Score         int            `json:"score"`
	FoundReport   ReportResponse `json:"found_report"`
	MissingReport ReportResponse `json:"missing_report"`
	Notified      bool           `json:"notified"`
	CreatedAt     time.Time      `json:"created_at"`
}

type GetMatchesQuery struct {
	MinScore int `form:"min_score" binding:"omitempty,min=0,max=100"`
	Page     int `form:"page"      binding:"omitempty,min=1"`
	Limit    int `form:"limit"     binding:"omitempty,min=1,max=100"`
}

type MatchListData struct {
	Matches []MatchResponse `json:"matches"`
	Meta    Pagination      `json:"meta"`
}
