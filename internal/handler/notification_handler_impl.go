package handler

import (
	"net/http"
	"temukan-api/internal/dto"
	"temukan-api/internal/exception"
	"temukan-api/internal/helper"
	"temukan-api/internal/usecase"

	"github.com/gin-gonic/gin"
)

type NotificationHandlerImpl struct {
	usecase usecase.NotificationUsecase
}

func NewNotificationHandlerImpl(usecase usecase.NotificationUsecase) *NotificationHandlerImpl {
	return &NotificationHandlerImpl{usecase: usecase}
}

// GET /notifications
func (h *NotificationHandlerImpl) GetAll(ctx *gin.Context) {
	userID := mustGetUserID(ctx)

	var query dto.GetNotificationsQuery
	if err := ctx.ShouldBindQuery(&query); err != nil {
		exception.ErrorHandler(ctx, err)
		return
	}

	result, err := h.usecase.GetAll(ctx, userID, query)
	if err != nil {
		exception.ErrorHandler(ctx, err)
		return
	}

	helper.WriteToResponseBody(ctx, http.StatusOK, dto.WebResponse{
		Status: "OK",
		Data:   result,
	})
}

// PATCH /notifications/read-all
// PENTING: route ini harus didaftarkan SEBELUM /notifications/:id/read
func (h *NotificationHandlerImpl) MarkAllAsRead(ctx *gin.Context) {
	userID := mustGetUserID(ctx)

	if err := h.usecase.MarkAllAsRead(ctx, userID); err != nil {
		exception.ErrorHandler(ctx, err)
		return
	}

	helper.WriteToResponseBody(ctx, http.StatusOK, dto.WebResponse{
		Status:  "OK",
		Message: "Semua notifikasi telah ditandai dibaca",
	})
}

// PATCH /notifications/:id/read
func (h *NotificationHandlerImpl) MarkAsRead(ctx *gin.Context) {
	userID := mustGetUserID(ctx)

	id, err := parseUUID(ctx, "id")
	if err != nil {
		exception.ErrorHandler(ctx, exception.NewBadRequestError("id tidak valid"))
		return
	}

	if err := h.usecase.MarkAsRead(ctx, id, userID); err != nil {
		exception.ErrorHandler(ctx, err)
		return
	}

	helper.WriteToResponseBody(ctx, http.StatusOK, dto.WebResponse{
		Status:  "OK",
		Message: "Notifikasi ditandai sudah dibaca",
	})
}
