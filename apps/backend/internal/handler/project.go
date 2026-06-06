package handler

import (
	"errors"
	"net/http"

	"github.com/Richard-OOO/E-director/apps/backend/internal/domain"
	"github.com/Richard-OOO/E-director/apps/backend/internal/response"
	"github.com/Richard-OOO/E-director/apps/backend/internal/service"
	"github.com/gin-gonic/gin"
)

type ProjectHandler struct {
	auth    *AuthHandler
	service *service.GenerationService
}

type createProjectRequest struct {
	Title      string `json:"title"`
	Language   string `json:"language"`
	SourceType string `json:"source_type"`
	Content    string `json:"content"`
}

func NewProjectHandler(auth *AuthHandler, generationService *service.GenerationService) *ProjectHandler {
	return &ProjectHandler{auth: auth, service: generationService}
}

func (h *ProjectHandler) Create(c *gin.Context) {
	user, ok := h.currentUser(c)
	if !ok {
		return
	}
	var payload createProjectRequest
	if err := c.ShouldBindJSON(&payload); err != nil {
		response.Error(c, http.StatusBadRequest, 40001, "invalid json body")
		return
	}
	result, err := h.service.CreateProjectAndStart(c.Request.Context(), service.CreateProjectInput{
		UserID:     user.ID,
		Title:      payload.Title,
		Language:   payload.Language,
		SourceType: payload.SourceType,
		Content:    payload.Content,
	})
	if err != nil {
		h.writeError(c, err)
		return
	}
	response.Created(c, gin.H{"project_id": result.ProjectID, "job_id": result.JobID, "status": result.Status})
}

func (h *ProjectHandler) List(c *gin.Context) {
	user, ok := h.currentUser(c)
	if !ok {
		return
	}
	items, err := h.service.ListProjects(c.Request.Context(), user.ID)
	if err != nil {
		h.writeError(c, err)
		return
	}
	response.OK(c, gin.H{"projects": items})
}

func (h *ProjectHandler) Detail(c *gin.Context) {
	user, ok := h.currentUser(c)
	if !ok {
		return
	}
	projectID := c.Param("project_id")
	snapshot, err := h.service.GetProject(c.Request.Context(), user.ID, projectID)
	if err != nil {
		h.writeError(c, err)
		return
	}
	response.OK(c, snapshot)
}

func (h *ProjectHandler) currentUser(c *gin.Context) (domain.User, bool) {
	if h.auth == nil || h.auth.service == nil {
		response.Error(c, http.StatusServiceUnavailable, 50001, "storage unavailable")
		return domain.User{}, false
	}
	token := h.auth.sessionToken(c)
	user, err := h.auth.service.Me(c.Request.Context(), token)
	if err != nil {
		h.auth.writeServiceError(c, err)
		return domain.User{}, false
	}
	return user, true
}

func (h *ProjectHandler) writeError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, service.ErrInvalidInput):
		response.Error(c, http.StatusBadRequest, 40001, "invalid input")
	case errors.Is(err, service.ErrUnauthorized):
		response.Error(c, http.StatusUnauthorized, 40100, "unauthorized")
	case errors.Is(err, service.ErrStorageUnavailable):
		response.Error(c, http.StatusServiceUnavailable, 50001, "storage unavailable")
	default:
		response.Error(c, http.StatusInternalServerError, 50000, "internal server error")
	}
}
