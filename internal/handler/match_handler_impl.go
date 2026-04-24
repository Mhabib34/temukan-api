package handler

import (
	"net/http"
	"temukan-api/internal/dto"
	"temukan-api/internal/exception"
	"temukan-api/internal/helper"
	"temukan-api/internal/usecase"

	"github.com/gin-gonic/gin"
)

type MatchHandlerImpl struct {
	usecase usecase.MatchUsecase
}

func NewMatchHandlerImpl(usecase usecase.MatchUsecase) *MatchHandlerImpl {
	return &MatchHandlerImpl{usecase: usecase}
}

// GET /matches
func (h *MatchHandlerImpl) GetAll(ctx *gin.Context) {
	userID := mustGetUserID(ctx)

	var query dto.GetMatchesQuery
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

// GET /matches/:id
func (h *MatchHandlerImpl) GetByID(ctx *gin.Context) {
	userID := mustGetUserID(ctx)

	id, err := parseUUID(ctx, "id")
	if err != nil {
		exception.ErrorHandler(ctx, exception.NewBadRequestError("id tidak valid"))
		return
	}

	result, err := h.usecase.GetByID(ctx, id, userID)
	if err != nil {
		exception.ErrorHandler(ctx, err)
		return
	}

	helper.WriteToResponseBody(ctx, http.StatusOK, dto.WebResponse{
		Status: "OK",
		Data:   result,
	})
}
