package handler

import (
	"net/http"
	"temukan-api/internal/dto"
	"temukan-api/internal/exception"
	"temukan-api/internal/helper"
	"temukan-api/internal/usecase"

	"github.com/gin-gonic/gin"
)

type StatsHandlerImpl struct {
	usecase usecase.StatsUsecase
}

func NewStatsHandlerImpl(usecase usecase.StatsUsecase) *StatsHandlerImpl {
	return &StatsHandlerImpl{usecase}
}

// GET /api/v1/stats
func (s *StatsHandlerImpl) GetStats(ctx *gin.Context) {
	result, err := s.usecase.GetStats(ctx)
	if err != nil {
		exception.ErrorHandler(ctx, err)
		return
	}

	helper.WriteToResponseBody(ctx, http.StatusOK, dto.WebResponse{
		Status: "OK",
		Data:   result,
	})
}