package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/Richard-OOO/E-director/apps/backend/internal/domain"
	"github.com/Richard-OOO/E-director/apps/backend/internal/response"
	"github.com/Richard-OOO/E-director/apps/backend/internal/service"
	"github.com/gin-gonic/gin"
)

type ProjectStreamHandler struct {
	auth   *AuthHandler
	events *service.GenerationEventBus
}

func NewProjectStreamHandler(auth *AuthHandler, events *service.GenerationEventBus) *ProjectStreamHandler {
	return &ProjectStreamHandler{auth: auth, events: events}
}

func (h *ProjectStreamHandler) Stream(c *gin.Context) {
	user, ok := h.currentUser(c)
	if !ok {
		return
	}
	_ = user
	if h.events == nil {
		response.Error(c, http.StatusServiceUnavailable, 50001, "stream unavailable")
		return
	}
	projectID := c.Param("project_id")
	jobID := c.Query("job_id")
	if projectID == "" || jobID == "" {
		response.Error(c, http.StatusBadRequest, 40001, "missing project_id or job_id")
		return
	}

	c.Header("Content-Type", "text/event-stream")
	c.Header("Cache-Control", "no-cache")
	c.Header("Connection", "keep-alive")
	c.Header("X-Accel-Buffering", "no")
	c.Status(http.StatusOK)

	eventsCh, cancel := h.events.Subscribe(projectID, jobID)
	defer cancel()

	flusher, ok := c.Writer.(http.Flusher)
	if !ok {
		response.Error(c, http.StatusInternalServerError, 50000, "streaming unsupported")
		return
	}

	ping := time.NewTicker(20 * time.Second)
	defer ping.Stop()

	writeEvent := func(eventName string, payload any) error {
		if _, err := fmt.Fprintf(c.Writer, "event: %s\n", eventName); err != nil {
			return err
		}
		data, err := json.Marshal(payload)
		if err != nil {
			return err
		}
		if _, err := fmt.Fprintf(c.Writer, "data: %s\n\n", data); err != nil {
			return err
		}
		flusher.Flush()
		return nil
	}

	for {
		select {
		case <-c.Request.Context().Done():
			return
		case <-ping.C:
			if _, err := fmt.Fprint(c.Writer, ": ping\n\n"); err != nil {
				return
			}
			flusher.Flush()
		case event, ok := <-eventsCh:
			if !ok {
				return
			}
			if err := writeEvent(string(event.EventType), event); err != nil {
				return
			}
			if event.EventType == domain.GenerationEventCompleted || event.EventType == domain.GenerationEventFailed {
				return
			}
		}
	}
}

func (h *ProjectStreamHandler) currentUser(c *gin.Context) (domain.User, bool) {
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
