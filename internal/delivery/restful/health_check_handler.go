package restful

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type HealthCheckResponse struct {
	Status string `json:"status"`
}

func (s *Server) HealthCheckHandler(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, HealthCheckResponse{
		Status: "ok",
	})
}
