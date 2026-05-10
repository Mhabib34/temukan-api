package usecase

import (
	"context"
	"temukan-api/internal/dto"
)

type StatsUsecase interface {
	GetStats(ctx context.Context) (*dto.StatsData, error)
}