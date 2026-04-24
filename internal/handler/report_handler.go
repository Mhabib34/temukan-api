package handler

import "github.com/gin-gonic/gin"

type ReportHandler interface {
	Create(ctx *gin.Context)
	GetAll(ctx *gin.Context)
	GetByID(ctx *gin.Context)
	GetMyReports(ctx *gin.Context)
	Update(ctx *gin.Context)
	Delete(ctx *gin.Context)
	UploadPhoto(ctx *gin.Context)
	GetMapPins(ctx *gin.Context)
}
