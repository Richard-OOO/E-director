package handler

import (
	"github.com/Richard-OOO/E-director/apps/gateway/internal/response"
	"github.com/gin-gonic/gin"
)

type SystemHandler struct{}

func NewSystemHandler() *SystemHandler {
	return &SystemHandler{}
}

func (h *SystemHandler) Health(c *gin.Context) {
	response.OK(c, map[string]string{"status": "healthy"})
}

func (h *SystemHandler) Version(c *gin.Context) {
	response.OK(c, map[string]string{"version": "v1"})
}
