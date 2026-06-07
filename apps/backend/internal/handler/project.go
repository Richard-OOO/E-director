package handler

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/Richard-OOO/E-director/apps/backend/internal/domain"
	"github.com/Richard-OOO/E-director/apps/backend/internal/response"
	"github.com/Richard-OOO/E-director/apps/backend/internal/service"
	"github.com/gin-gonic/gin"
)

type ProjectHandler struct {
	auth           *AuthHandler
	service        *service.GenerationService
	extractor      service.DocumentTextExtractor
	importMaxBytes int64
}

type createProjectRequest struct {
	Title      string `json:"title"`
	Language   string `json:"language"`
	SourceType string `json:"source_type"`
	Content    string `json:"content"`
}

type updateSceneYAMLRequest struct {
	YAML string `json:"yaml"`
}

func NewProjectHandler(auth *AuthHandler, generationService *service.GenerationService, extractor service.DocumentTextExtractor, importMaxBytes int64) *ProjectHandler {
	return &ProjectHandler{auth: auth, service: generationService, extractor: extractor, importMaxBytes: importMaxBytes}
}

func (h *ProjectHandler) Create(c *gin.Context) {
	user, ok := h.currentUser(c)
	if !ok {
		return
	}

	payload, err := h.createProjectPayload(c)
	if err != nil {
		h.writeError(c, err)
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
		log.Printf("[e-director:project] create failed user_id=%s title=%q source_type=%q content_len=%d err=%v", user.ID, payload.Title, payload.SourceType, len(payload.Content), err)
		h.writeError(c, err)
		return
	}
	log.Printf("[e-director:project] create ok user_id=%s project_id=%s job_id=%s status=%s title=%q source_type=%q content_len=%d", user.ID, result.ProjectID, result.JobID, result.Status, payload.Title, payload.SourceType, len(payload.Content))
	response.Created(c, gin.H{"project_id": result.ProjectID, "job_id": result.JobID, "status": result.Status})
}

func (h *ProjectHandler) createProjectPayload(c *gin.Context) (createProjectRequest, error) {
	contentType := c.GetHeader("Content-Type")
	if strings.HasPrefix(contentType, "multipart/form-data") {
		return h.createProjectPayloadFromMultipart(c)
	}

	var payload createProjectRequest
	if err := c.ShouldBindJSON(&payload); err != nil {
		return createProjectRequest{}, service.ErrInvalidInput
	}
	return payload, nil
}

func (h *ProjectHandler) createProjectPayloadFromMultipart(c *gin.Context) (createProjectRequest, error) {
	if h.extractor == nil {
		return createProjectRequest{}, service.ErrStorageUnavailable
	}
	maxBytes := h.importMaxBytes
	if maxBytes <= 0 {
		maxBytes = 20 << 20
	}
	c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, maxBytes)
	if err := c.Request.ParseMultipartForm(maxBytes); err != nil {
		return createProjectRequest{}, service.ErrInvalidInput
	}

	file, header, err := c.Request.FormFile("file")
	if err != nil {
		return createProjectRequest{}, service.ErrInvalidInput
	}
	defer file.Close()

	filename := strings.TrimSpace(header.Filename)
	sourceType := strings.TrimSpace(c.Request.FormValue("source_type"))
	if sourceType == "" {
		sourceType = strings.TrimPrefix(strings.ToLower(filepath.Ext(filename)), ".")
	}
	text, err := h.extractor.Extract(c.Request.Context(), service.DocumentExtractInput{
		SourceType:  sourceType,
		Filename:    filename,
		ContentType: header.Header.Get("Content-Type"),
		Reader:      file,
		MaxBytes:    maxBytes,
	})
	if err != nil {
		return createProjectRequest{}, err
	}

	title := strings.TrimSpace(c.Request.FormValue("title"))
	if title == "" {
		title = strings.TrimSuffix(filename, filepath.Ext(filename))
	}
	return createProjectRequest{
		Title:      title,
		Language:   c.Request.FormValue("language"),
		SourceType: sourceType,
		Content:    text,
	}, nil
}

func (h *ProjectHandler) List(c *gin.Context) {
	user, ok := h.currentUser(c)
	if !ok {
		return
	}
	items, err := h.service.ListProjects(c.Request.Context(), user.ID)
	if err != nil {
		log.Printf("[e-director:project] list failed user_id=%s err=%v", user.ID, err)
		h.writeError(c, err)
		return
	}
	log.Printf("[e-director:project] list ok user_id=%s projects=%d", user.ID, len(items))
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
		log.Printf("[e-director:project] detail failed user_id=%s project_id=%s err=%v", user.ID, projectID, err)
		h.writeError(c, err)
		return
	}
	response.OK(c, snapshot)
}

func (h *ProjectHandler) UpdateSceneYAML(c *gin.Context) {
	user, ok := h.currentUser(c)
	if !ok {
		return
	}
	projectID := c.Param("project_id")
	sceneID := c.Param("scene_id")
	var payload updateSceneYAMLRequest
	if err := c.ShouldBindJSON(&payload); err != nil {
		h.writeError(c, service.ErrInvalidInput)
		return
	}
	if err := h.service.UpdateSceneYAML(c.Request.Context(), service.UpdateSceneYAMLInput{UserID: user.ID, ProjectID: projectID, SceneID: sceneID, YAML: payload.YAML}); err != nil {
		log.Printf("[e-director:project] update scene yaml failed user_id=%s project_id=%s scene_id=%s yaml_len=%d err=%v", user.ID, projectID, sceneID, len(payload.YAML), err)
		h.writeError(c, err)
		return
	}
	log.Printf("[e-director:project] update scene yaml ok user_id=%s project_id=%s scene_id=%s yaml_len=%d", user.ID, projectID, sceneID, len(payload.YAML))
	response.OK(c, gin.H{"scene_id": sceneID})
}

func (h *ProjectHandler) Delete(c *gin.Context) {
	user, ok := h.currentUser(c)
	if !ok {
		return
	}
	projectID := c.Param("project_id")
	if err := h.service.DeleteProject(c.Request.Context(), user.ID, projectID); err != nil {
		log.Printf("[e-director:project] delete failed user_id=%s project_id=%s err=%v", user.ID, projectID, err)
		h.writeError(c, err)
		return
	}
	log.Printf("[e-director:project] delete ok user_id=%s project_id=%s", user.ID, projectID)
	response.OK(c, gin.H{"project_id": projectID})
}

func (h *ProjectHandler) ExportYAML(c *gin.Context) {
	user, ok := h.currentUser(c)
	if !ok {
		return
	}
	projectID := c.Param("project_id")
	result, err := h.service.ExportProjectYAML(c.Request.Context(), user.ID, projectID)
	if err != nil {
		log.Printf("[e-director:project] export yaml failed user_id=%s project_id=%s err=%v", user.ID, projectID, err)
		h.writeError(c, err)
		return
	}
	h.writeYAMLExport(c, result)
}

func (h *ProjectHandler) ExportCombinedYAML(c *gin.Context) {
	user, ok := h.currentUser(c)
	if !ok {
		return
	}
	projectID := c.Param("project_id")
	result, err := h.service.ExportCombinedProjectYAML(c.Request.Context(), user.ID, projectID)
	if err != nil {
		log.Printf("[e-director:project] export combined yaml failed user_id=%s project_id=%s err=%v", user.ID, projectID, err)
		h.writeError(c, err)
		return
	}
	h.writeYAMLExport(c, result)
}

func (h *ProjectHandler) writeYAMLExport(c *gin.Context, result service.ExportProjectYAMLResult) {
	c.Header("Content-Type", result.ContentType)
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%q", result.Filename))
	c.Header("X-Content-Type-Options", "nosniff")
	c.Data(http.StatusOK, result.ContentType, result.Data)
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
	case errors.Is(err, service.ErrUnsupportedDocument):
		response.Error(c, http.StatusBadRequest, 40002, "unsupported document type")
	case errors.Is(err, service.ErrEmptyDocumentText):
		response.Error(c, http.StatusBadRequest, 40003, "document contains no extractable text")
	case errors.Is(err, service.ErrUnauthorized):
		response.Error(c, http.StatusUnauthorized, 40100, "unauthorized")
	case errors.Is(err, service.ErrNotFound):
		response.Error(c, http.StatusNotFound, 40401, "project or scene not found")
	case errors.Is(err, service.ErrStorageUnavailable):
		response.Error(c, http.StatusServiceUnavailable, 50001, "storage unavailable")
	default:
		response.Error(c, http.StatusInternalServerError, 50000, "internal server error")
	}
}
