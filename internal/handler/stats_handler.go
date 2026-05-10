package handler

import "github.com/gin-gonic/gin"

type StatsHandler interface {
	GetStats(ctx *gin.Context)
}