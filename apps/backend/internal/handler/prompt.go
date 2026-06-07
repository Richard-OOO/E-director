package handler

import (
	"errors"
	"net/http"

	"github.com/Richard-OOO/E-director/apps/backend/internal/domain"
	"github.com/Richard-OOO/E-director/apps/backend/internal/response"
	"github.com/Richard-OOO/E-director/apps/backend/internal/service"
	"github.com/gin-gonic/gin"
)

type PromptHandler struct {
	auth    *AuthHandler
	service *service.PromptService
}

func NewPromptHandler(auth *AuthHandler, promptService *service.PromptService) *PromptHandler {
	return &PromptHandler{auth: auth, service: promptService}
}

type promptRequest struct {
	Name    string `json:"name"`
	Content string `json:"content"`
}

func (h *PromptHandler) List(c *gin.Context) {
	user, ok := h.currentUser(c)
	if !ok {
		return
	}
	prompts, err := h.service.ListPrompts(c.Request.Context(), user.ID)
	if err != nil {
		h.writeError(c, err)
		return
	}
	response.OK(c, gin.H{"prompts": prompts})
}

func (h *PromptHandler) Detail(c *gin.Context) {
	user, ok := h.currentUser(c)
	if !ok {
		return
	}
	promptID := c.Param("prompt_id")
	prompt, err := h.service.GetPrompt(c.Request.Context(), user.ID, promptID)
	if err != nil {
		h.writeError(c, err)
		return
	}
	response.OK(c, prompt)
}

func (h *PromptHandler) Create(c *gin.Context) {
	user, ok := h.currentUser(c)
	if !ok {
		return
	}
	payload, err := h.bindPromptRequest(c)
	if err != nil {
		h.writeError(c, err)
		return
	}
	prompt, err := h.service.CreatePrompt(c.Request.Context(), service.CreatePromptInput{UserID: user.ID, Name: payload.Name, Content: payload.Content})
	if err != nil {
		h.writeError(c, err)
		return
	}
	response.Created(c, prompt)
}

func (h *PromptHandler) Update(c *gin.Context) {
	user, ok := h.currentUser(c)
	if !ok {
		return
	}
	payload, err := h.bindPromptRequest(c)
	if err != nil {
		h.writeError(c, err)
		return
	}
	promptID := c.Param("prompt_id")
	prompt, err := h.service.UpdatePrompt(c.Request.Context(), service.UpdatePromptInput{UserID: user.ID, PromptID: promptID, Name: payload.Name, Content: payload.Content})
	if err != nil {
		h.writeError(c, err)
		return
	}
	response.OK(c, prompt)
}

func (h *PromptHandler) Delete(c *gin.Context) {
	user, ok := h.currentUser(c)
	if !ok {
		return
	}
	promptID := c.Param("prompt_id")
	if err := h.service.DeletePrompt(c.Request.Context(), user.ID, promptID); err != nil {
		h.writeError(c, err)
		return
	}
	response.OK(c, response.EmptyData())
}

func (h *PromptHandler) Reset(c *gin.Context) {
	user, ok := h.currentUser(c)
	if !ok {
		return
	}
	promptID := c.Param("prompt_id")
	prompt, err := h.service.ResetPrompt(c.Request.Context(), user.ID, promptID)
	if err != nil {
		h.writeError(c, err)
		return
	}
	response.OK(c, prompt)
}

func (h *PromptHandler) bindPromptRequest(c *gin.Context) (promptRequest, error) {
	var payload promptRequest
	if err := c.ShouldBindJSON(&payload); err != nil {
		return promptRequest{}, service.ErrInvalidInput
	}
	return payload, nil
}

func (h *PromptHandler) currentUser(c *gin.Context) (domain.User, bool) {
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

func (h *PromptHandler) writeError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, service.ErrInvalidInput):
		response.Error(c, http.StatusBadRequest, 40001, "invalid input")
	case errors.Is(err, service.ErrUnauthorized):
		response.Error(c, http.StatusUnauthorized, 40100, "unauthorized")
	case errors.Is(err, service.ErrNotFound):
		response.Error(c, http.StatusNotFound, 40401, "prompt not found")
	case errors.Is(err, service.ErrStorageUnavailable):
		response.Error(c, http.StatusServiceUnavailable, 50001, "storage unavailable")
	default:
		response.Error(c, http.StatusInternalServerError, 50000, "internal server error")
	}
}
