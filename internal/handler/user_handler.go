package handler

import "github.com/gin-gonic/gin"

type UserHandler interface {
	Create(ctx *gin.Context)
	Login(ctx *gin.Context)
	Logout(ctx *gin.Context)
	RefreshToken(ctx *gin.Context)
	Me(ctx *gin.Context)
}
