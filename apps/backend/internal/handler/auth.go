package handler

import (
	"errors"
	"net/http"
	"strings"

	"github.com/Richard-OOO/E-director/apps/backend/internal/config"
	"github.com/Richard-OOO/E-director/apps/backend/internal/domain"
	"github.com/Richard-OOO/E-director/apps/backend/internal/response"
	"github.com/Richard-OOO/E-director/apps/backend/internal/service"
	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	service *service.AuthService
	cfg     config.Config
}

type authRequest struct {
	Email            string `json:"email"`
	Password         string `json:"password"`
	VerificationCode string `json:"verification_code"`
	DisplayName      string `json:"display_name"`
	Name             string `json:"name"`
}

type authResponse struct {
	User userDTO `json:"user"`
}

type userDTO struct {
	ID          string `json:"id"`
	Email       string `json:"email"`
	DisplayName string `json:"display_name"`
}

func NewAuthHandler(authService *service.AuthService, cfg config.Config) *AuthHandler {
	return &AuthHandler{service: authService, cfg: cfg}
}

func (h *AuthHandler) SendCode(c *gin.Context) {
	var payload authRequest
	if err := c.ShouldBindJSON(&payload); err != nil {
		response.Error(c, http.StatusBadRequest, 40001, "invalid json body")
		return
	}
	result, err := h.service.SendCode(c.Request.Context(), payload.Email)
	if err != nil {
		h.writeServiceError(c, err)
		return
	}
	response.OK(c, result)
}

func (h *AuthHandler) Register(c *gin.Context) {
	input, ok := bindAuthInput(c)
	if !ok {
		return
	}
	result, err := h.service.Register(c.Request.Context(), input)
	if err != nil {
		h.writeServiceError(c, err)
		return
	}
	h.setSessionCookie(c, result.SessionToken)
	response.Created(c, toAuthResponse(result))
}

func (h *AuthHandler) Login(c *gin.Context) {
	input, ok := bindAuthInput(c)
	if !ok {
		return
	}
	result, err := h.service.Login(c.Request.Context(), input)
	if err != nil {
		h.writeServiceError(c, err)
		return
	}
	h.setSessionCookie(c, result.SessionToken)
	response.OK(c, toAuthResponse(result))
}

func (h *AuthHandler) Me(c *gin.Context) {
	token := h.sessionToken(c)
	user, err := h.service.Me(c.Request.Context(), token)
	if err != nil {
		h.writeServiceError(c, err)
		return
	}
	response.OK(c, toUserDTO(user))
}

func (h *AuthHandler) Logout(c *gin.Context) {
	token := h.sessionToken(c)
	if err := h.service.Logout(c.Request.Context(), token); err != nil {
		h.writeServiceError(c, err)
		return
	}
	h.clearSessionCookie(c)
	response.OK(c, response.EmptyData())
}

func bindAuthInput(c *gin.Context) (service.AuthInput, bool) {
	var payload authRequest
	if err := c.ShouldBindJSON(&payload); err != nil {
		response.Error(c, http.StatusBadRequest, 40001, "invalid json body")
		return service.AuthInput{}, false
	}
	displayName := strings.TrimSpace(payload.DisplayName)
	if displayName == "" {
		displayName = strings.TrimSpace(payload.Name)
	}
	return service.AuthInput{Email: payload.Email, Password: payload.Password, VerificationCode: payload.VerificationCode, DisplayName: displayName}, true
}

func (h *AuthHandler) writeServiceError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, service.ErrInvalidInput):
		response.Error(c, http.StatusBadRequest, 40001, "invalid input")
	case errors.Is(err, service.ErrEmailRegistered):
		response.Error(c, http.StatusConflict, 40002, "email already registered")
	case errors.Is(err, service.ErrInvalidCredentials):
		response.Error(c, http.StatusUnauthorized, 40101, "invalid email or password")
	case errors.Is(err, service.ErrCodeExpired):
		response.Error(c, http.StatusUnauthorized, 40102, "invalid or expired verification code")
	case errors.Is(err, service.ErrTooManyAttempts):
		response.Error(c, http.StatusTooManyRequests, 42902, "too many verification attempts")
	case errors.Is(err, service.ErrCodeCooldown):
		response.Error(c, http.StatusTooManyRequests, 42901, "verification code sent too frequently")
	case errors.Is(err, service.ErrUnauthorized):
		response.Error(c, http.StatusUnauthorized, 40100, "unauthorized")
	case errors.Is(err, service.ErrStorageUnavailable):
		response.Error(c, http.StatusServiceUnavailable, 50001, "storage unavailable")
	case errors.Is(err, service.ErrMailUnavailable):
		response.Error(c, http.StatusServiceUnavailable, 50002, "mail service unavailable")
	default:
		response.Error(c, http.StatusInternalServerError, 50000, "internal server error")
	}
}

func (h *AuthHandler) setSessionCookie(c *gin.Context, token string) {
	maxAge := int(h.cfg.SessionTTL.Seconds())
	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie(h.cfg.CookieName, token, maxAge, "/", "", h.cfg.CookieSecure, true)
}

func (h *AuthHandler) clearSessionCookie(c *gin.Context) {
	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie(h.cfg.CookieName, "", -1, "/", "", h.cfg.CookieSecure, true)
}

func (h *AuthHandler) sessionToken(c *gin.Context) string {
	if token, err := c.Cookie(h.cfg.CookieName); err == nil && strings.TrimSpace(token) != "" {
		return strings.TrimSpace(token)
	}
	return bearerToken(c.GetHeader("Authorization"))
}

func toAuthResponse(result domain.AuthResult) authResponse {
	return authResponse{User: toUserDTO(result.User)}
}

func toUserDTO(user domain.User) userDTO {
	return userDTO{ID: user.ID, Email: user.Email, DisplayName: user.DisplayName}
}

func bearerToken(header string) string {
	const prefix = "Bearer "
	if strings.HasPrefix(header, prefix) {
		return strings.TrimSpace(strings.TrimPrefix(header, prefix))
	}
	return strings.TrimSpace(header)
}
