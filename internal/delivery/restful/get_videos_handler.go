package restful

import (
	db "github.com/dinhtatuanlinh/video/db/sqlc"
	apierr "github.com/dinhtatuanlinh/video/internal/delivery/restful/error"
	internalError "github.com/dinhtatuanlinh/video/internal/error"
	"github.com/dinhtatuanlinh/video/internal/usecase/video"
	"github.com/dinhtatuanlinh/video/util"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"net/http"
	"time"
)

type GetVideoRequest struct {
	Limit  int32 `form:"limit"`
	Offset int32 `form:"offset"`
}

type GetVideoResponse struct {
	Videos []VideoResponse `json:"videos"`
}
type VideoResponse struct {
	VideoID           int64     `json:"video_id"`
	VideoCategoryName string    `json:"video_category_name"`
	Name              string    `json:"name"`
	Url               string    `json:"url"`
	CreatedAt         time.Time `json:"created_at"`
}

func (s *Server) GetVideosHandler(ctx *gin.Context) {
	log.Info().Msg("CreateVideoCategoryHandler called")
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
	videos, err := s.useCase.UseCaseVideo.GetVideos(ctx, request)
	if err != nil {
		log.Error().Err(err).Msg("Failed to CreateVideoCategory")
		ctx.JSON(http.StatusInternalServerError, apierr.ErrorResponse(err))
		return
	}
	videosRes := util.Map(func(video db.Video) VideoResponse {
		return VideoResponse{
			VideoID:           video.VideoID,
			VideoCategoryName: video.VideoCategoryName,
			Name:              video.Name,
			Url:               video.Url,
		}
	}, videos)
	ctx.JSON(http.StatusOK, GetVideoResponse{
		videosRes,
	})
}
