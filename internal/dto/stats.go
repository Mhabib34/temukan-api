package dto

type StatsData struct {
	ActiveReports   int64 `json:"active_reports"`
	TotalVolunteers int64 `json:"total_volunteers"`
	ResolvedLast24h int64 `json:"resolved_last_24h"`
	UniqueCities    int64 `json:"unique_cities"`
	TotalResolved   int64 `json:"total_resolved"`
}