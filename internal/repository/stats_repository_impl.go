package repository

import (
	"context"
	"time"

	"gorm.io/gorm"
)

type StatsRepositoryImpl struct {
	DB *gorm.DB
}

func NewStatsRepository(db *gorm.DB) StatsRepository {
	return &StatsRepositoryImpl{DB: db}
}

func (s *StatsRepositoryImpl) CountActiveReports(ctx context.Context) (int64, error) {
	var count int64
	err := s.DB.WithContext(ctx).
		Table("reports").
		Where("status = ?", "active").
		Count(&count).Error
	return count, err
}

func (s *StatsRepositoryImpl) CountVolunteers(ctx context.Context) (int64, error) {
	var count int64
	err := s.DB.WithContext(ctx).
		Table("users").
		Where("role = ?", "volunteer").
		Count(&count).Error
	return count, err
}

func (s *StatsRepositoryImpl) CountResolvedLast24h(ctx context.Context) (int64, error) {
	var count int64
	since := time.Now().Add(-24 * time.Hour)
	err := s.DB.WithContext(ctx).
		Table("reports").
		Where("status = ? AND updated_at >= ?", "resolved", since).
		Count(&count).Error
	return count, err
}

func (s *StatsRepositoryImpl) CountUniqueCities(ctx context.Context) (int64, error) {
	var count int64
	err := s.DB.WithContext(ctx).
		Table("reports").
		Where("status = ?", "active").
		Distinct("city").
		Count(&count).Error
	return count, err
}

func (s *StatsRepositoryImpl) CountTotalResolved(ctx context.Context) (int64, error) {
	var count int64
	err := s.DB.WithContext(ctx).
		Table("reports").
		Where("status = ?", "resolved").
		Count(&count).Error
	return count, err
}