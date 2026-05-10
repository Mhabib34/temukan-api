package repository

import "context"

type StatsRepository interface {
	CountActiveReports(ctx context.Context) (int64, error)
	CountVolunteers(ctx context.Context) (int64, error)
	CountResolvedLast24h(ctx context.Context) (int64, error)
	CountUniqueCities(ctx context.Context) (int64, error)
}