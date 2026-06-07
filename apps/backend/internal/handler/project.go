package handler

import (
	"errors"
	"log"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/Richard-OOO/E-director/apps/backend/internal/domain"
	mysqlmodels "github.com/Richard-OOO/E-director/apps/backend/internal/models/mysql"
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

func NewProjectHandler(auth *AuthHandler, generationService *service.GenerationService, extractor service.DocumentTextExtractor, importMaxBytes int64) *ProjectHandler {
	return &ProjectHandler{auth: auth, service: generationService, extractor: extractor, importMaxBytes: importMaxBytes}
}

func debugProjectSnapshot(snapshot mysqlmodels.ProjectSnapshot) {
	firstSceneEditableYAMLLen := 0
	firstSceneGeneratedYAMLLen := 0
	firstSceneDesignReasonYAMLLen := 0
	if len(snapshot.Scenes) > 0 {
		firstSceneEditableYAMLLen = len(snapshot.Scenes[0].EditableYAML)
		firstSceneGeneratedYAMLLen = len(snapshot.Scenes[0].GeneratedYAML)
		firstSceneDesignReasonYAMLLen = len(snapshot.Scenes[0].DesignReasonYAML)
	}
	log.Printf("[e-director:project] detail snapshot project_id=%s job_id=%s chapters=%d scenes=%d first_scene_editable_yaml_len=%d first_scene_generated_yaml_len=%d first_scene_design_reason_yaml_len=%d",
		snapshot.Project.ID,
		snapshot.Job.ID,
		len(snapshot.Chapters),
		len(snapshot.Scenes),
		firstSceneEditableYAMLLen,
		firstSceneGeneratedYAMLLen,
		firstSceneDesignReasonYAMLLen,
	)
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
	debugProjectSnapshot(snapshot)
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
	case errors.Is(err, service.ErrUnsupportedDocument):
		response.Error(c, http.StatusBadRequest, 40002, "unsupported document type")
	case errors.Is(err, service.ErrEmptyDocumentText):
		response.Error(c, http.StatusBadRequest, 40003, "document contains no extractable text")
	case errors.Is(err, service.ErrUnauthorized):
		response.Error(c, http.StatusUnauthorized, 40100, "unauthorized")
	case errors.Is(err, service.ErrStorageUnavailable):
		response.Error(c, http.StatusServiceUnavailable, 50001, "storage unavailable")
	default:
		response.Error(c, http.StatusInternalServerError, 50000, "internal server error")
	}
}
