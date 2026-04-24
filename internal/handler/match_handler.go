package handler

import "github.com/gin-gonic/gin"

type MatchHandler interface {
	GetAll(ctx *gin.Context)
	GetByID(ctx *gin.Context)
}
