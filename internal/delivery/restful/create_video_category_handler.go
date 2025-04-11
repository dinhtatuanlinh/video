package restful

import (
	apierr "github.com/dinhtatuanlinh/video/internal/delivery/restful/error"
	internalError "github.com/dinhtatuanlinh/video/internal/error"
	"github.com/dinhtatuanlinh/video/internal/usecase/video"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"net/http"
)

type CreateVideoCategoryRequest struct {
	CategoryName       string `json:"category_name"`
	CategoryParentName string `json:"category_parent_name"`
}

func (s *Server) CreateVideoCategoryHandler(ctx *gin.Context) {
	log.Info().Msg("CreateVideoCategoryHandler called")
	var req CreateVideoCategoryRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Error().Err(err).Msg("Failed to bind JSON body")
		if appErr, ok := err.(*internalError.AppError); ok {
			ctx.JSON(appErr.HTTPCode, gin.H{
				"error": appErr.Error(), // includes internal error
			})
			return
		}
		ctx.JSON(http.StatusInternalServerError, apierr.ErrorResponse(err))
		return
	}

	request := &video.CreateVideoCategoryModel{
		CategoryName:       req.CategoryName,
		CategoryParentName: req.CategoryParentName,
	}
	err := s.useCase.UseCaseVideo.CreateVideoCategory(ctx, request)
	if err != nil {
		log.Error().Err(err).Msg("Failed to CreateVideoCategory")
		ctx.JSON(http.StatusInternalServerError, apierr.ErrorResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, struct{}{})
}
