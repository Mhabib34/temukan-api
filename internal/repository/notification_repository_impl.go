package repository

import (
	"context"
	"temukan-api/internal/dto"
	"temukan-api/internal/model"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type NotificationRepositoryImpl struct {
	DB *gorm.DB
}

func NewNotificationRepository(db *gorm.DB) NotificationRepository {
	return &NotificationRepositoryImpl{DB: db}
}

func (r *NotificationRepositoryImpl) Create(ctx context.Context, notification *model.Notification) (*model.Notification, error) {
	err := r.DB.WithContext(ctx).Create(notification).Error
	if err != nil {
		return nil, err
	}
	return notification, nil
}

func (r *NotificationRepositoryImpl) FindByUserID(ctx context.Context, userID uuid.UUID, query dto.GetNotificationsQuery) ([]model.Notification, int64, error) {
	var notifications []model.Notification
	var total int64

	db := r.DB.WithContext(ctx).
		Model(&model.Notification{}).
		Where("user_id = ?", userID)

	if query.IsRead != nil {
		db = db.Where("is_read = ?", *query.IsRead)
	}

	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	page := query.Page
	if page < 1 {
		page = 1
	}
	limit := query.Limit
	if limit < 1 {
		limit = 10
	}
	offset := (page - 1) * limit

	err := db.Order("created_at DESC").Offset(offset).Limit(limit).Find(&notifications).Error
	return notifications, total, err
}

func (r *NotificationRepositoryImpl) CountUnread(ctx context.Context, userID uuid.UUID) (int64, error) {
	var count int64
	err := r.DB.WithContext(ctx).
		Model(&model.Notification{}).
		Where("user_id = ? AND is_read = false", userID).
		Count(&count).Error
	return count, err
}

func (r *NotificationRepositoryImpl) MarkAsRead(ctx context.Context, id uuid.UUID, userID uuid.UUID) error {
	result := r.DB.WithContext(ctx).
		Model(&model.Notification{}).
		Where("id = ? AND user_id = ?", id, userID).
		Update("is_read", true)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func (r *NotificationRepositoryImpl) MarkAllAsRead(ctx context.Context, userID uuid.UUID) error {
	return r.DB.WithContext(ctx).
		Model(&model.Notification{}).
		Where("user_id = ? AND is_read = false", userID).
		Update("is_read", true).Error
}
