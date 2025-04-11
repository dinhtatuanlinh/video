package restful

import (
	apierr "github.com/dinhtatuanlinh/video/internal/delivery/restful/error"
	internalError "github.com/dinhtatuanlinh/video/internal/error"
	"github.com/dinhtatuanlinh/video/internal/usecase/video"
	"github.com/dinhtatuanlinh/video/util"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"net/http"
)

type CreateVideoRequest struct {
	Videos []VideoRequest `json:"videos"`
}
type VideoRequest struct {
	CategoryName string `json:"category_name"`
	Name         string `json:"name"`
	InputPath    string `json:"input_path"`
}

func (s *Server) CreateVideoHandler(ctx *gin.Context) {
	log.Info().Msg("CreateVideoHandler called")
	var req CreateVideoRequest
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
	videos := util.Map(func(v VideoRequest) video.VideoModel {
		return video.VideoModel{
			CategoryName: v.CategoryName,
			Name:         v.Name,
			InputPath:    v.InputPath,
		}
	}, req.Videos)
	request := &video.CreateVideoModel{
		Videos: videos,
	}
	err := s.useCase.UseCaseVideo.CreateVideo(ctx, request)
	if err != nil {
		log.Error().Err(err).Msg("Failed to CreateOperator")
		ctx.JSON(http.StatusInternalServerError, apierr.ErrorResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, struct{}{})
}
