package restful

import (
	apierr "github.com/dinhtatuanlinh/video/internal/delivery/restful/error"
	internalError "github.com/dinhtatuanlinh/video/internal/error"
	"github.com/dinhtatuanlinh/video/internal/usecase/video"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"net/http"
)

type DeleteVideoRequest struct {
	ID int64 `uri:"id" binding:"required"`
}

func (s *Server) DeleteVideoHandler(ctx *gin.Context) {
	log.Info().Msg("CreateVideoHandler called")
	var req DeleteVideoRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
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

	request := &video.DeleteVideosModel{
		ID: req.ID,
	}
	err := s.useCase.UseCaseVideo.DeleteVideos(ctx, request)
	if err != nil {
		log.Error().Err(err).Msg("Failed to DeleteVideos")
		ctx.JSON(http.StatusInternalServerError, apierr.ErrorResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, struct{}{})
}
