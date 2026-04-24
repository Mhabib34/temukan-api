package handler

import "github.com/gin-gonic/gin"

type NotificationHandler interface {
	GetAll(ctx *gin.Context)
	MarkAsRead(ctx *gin.Context)
	MarkAllAsRead(ctx *gin.Context)
}
