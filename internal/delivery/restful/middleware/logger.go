package middleware

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"net/http"
	"time"
)

func HttpLogger() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		startTime := time.Now()
		ctx.Next()
		duration := time.Since(startTime)

		logger := log.Info()
		if ctx.Writer.Status() != http.StatusOK {
			logger = log.Error()
			err := ctx.Errors
			fmt.Print(err)
			for _, ginErr := range ctx.Errors {
				logger.Str("error", ginErr.Error())
			}
		}

		logger.Str("protocol", "http").
			Str("method", ctx.Request.Method).
			Str("path", ctx.Request.URL.Path).
			Str("payload", ctx.Request.PostForm.Encode()).
			Int("status_code", ctx.Writer.Status()).
			Str("status_text", http.StatusText(ctx.Writer.Status())).
			Dur("duration", duration).
			Str("start_time", startTime.String()).
			Msg("received a http request")
	}
}
