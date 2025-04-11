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

type DownloadVideoRequest struct {
	Urls []UrlRequest `json:"urls"`
}
type UrlRequest struct {
	CategoryName string `json:"category_name"`
	Url          string `json:"url"`
}

func (s *Server) DownloadVideoHandler(ctx *gin.Context) {
	log.Info().Msg("DownloadVideoHandler called")
	var req DownloadVideoRequest
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
	urls := util.Map(func(url UrlRequest) video.UrlModel {
		return video.UrlModel{
			CategoryName: url.CategoryName,
			Url:          url.Url,
		}
	}, req.Urls)
	request := &video.DownloadVideoModel{
		Urls: urls,
	}
	err := s.useCase.UseCaseVideo.DownloadVideo(ctx, request)
	if err != nil {
		log.Error().Err(err).Msg("Failed to CreateOperator")
		ctx.JSON(http.StatusInternalServerError, apierr.ErrorResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, struct{}{})
}
