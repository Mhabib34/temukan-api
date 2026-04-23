package handler

import "github.com/gin-gonic/gin"

type UserHandler interface {
	Create(ctx *gin.Context)
}
