package restful

import (
	apierr "github.com/dinhtatuanlinh/video/internal/delivery/restful/error"
	internalError "github.com/dinhtatuanlinh/video/internal/error"
	"github.com/dinhtatuanlinh/video/internal/usecase/video"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"net/http"
)

type GetFileRequest struct {
	FolderName string `json:"folder_name"`
}

func (s *Server) GetFilesHandler(ctx *gin.Context) {
	log.Info().Msg("GetFilesHandler called")
	var req GetVideoRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
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

	request := &video.GetVideosModel{
		Limit:  req.Limit,
		Offset: req.Offset,
	}
	_, err := s.useCase.UseCaseVideo.GetVideos(ctx, request)
	if err != nil {
		log.Error().Err(err).Msg("Failed to GetVideos")
		ctx.JSON(http.StatusInternalServerError, apierr.ErrorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, struct{}{})
}
