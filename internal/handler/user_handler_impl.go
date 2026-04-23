package handler

import (
	"net/http"
	"temukan-api/internal/dto"
	"temukan-api/internal/exception"
	"temukan-api/internal/helper"
	"temukan-api/internal/usecase"

	"github.com/gin-gonic/gin"
)

type UserHandlerImpl struct {
	usecase usecase.UserUsecase
}

func NewUserHandlerImpl(usecase usecase.UserUsecase) *UserHandlerImpl {
	return &UserHandlerImpl{usecase}
}

func (u *UserHandlerImpl) Create(ctx *gin.Context) {
	var request dto.RegisterRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		exception.ErrorHandler(ctx, err)
		return
	}

	result, err := u.usecase.Create(ctx, &request)
	if err != nil {
		exception.ErrorHandler(ctx, err)
		return
	}

	helper.WriteToResponseBody(ctx, http.StatusCreated, dto.WebResponse{
		Status:  "OK",
		Message: "User created successfully",
		Data:    result,
	})
}
