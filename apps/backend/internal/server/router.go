package server

import (
	"github.com/Richard-OOO/E-director/apps/backend/internal/config"
	"github.com/Richard-OOO/E-director/apps/backend/internal/handler"
	"github.com/Richard-OOO/E-director/apps/backend/internal/middleware"
	"github.com/gin-gonic/gin"
)

type Handlers struct {
	System        *handler.SystemHandler
	Auth          *handler.AuthHandler
	Project       *handler.ProjectHandler
	ProjectStream *handler.ProjectStreamHandler
}

func NewRouter(cfg config.Config, handlers Handlers) *gin.Engine {
	router := gin.New()
	router.Use(gin.Logger(), gin.Recovery(), middleware.CORS(cfg.FrontendOrigins))

	router.GET("/health", handlers.System.Health)

	api := router.Group("/api/v1")
	api.GET("/version", handlers.System.Version)

	auth := api.Group("/auth")
	auth.POST("/send-code", handlers.Auth.SendCode)
	auth.POST("/register", handlers.Auth.Register)
	auth.POST("/login", handlers.Auth.Login)
	auth.POST("/logout", handlers.Auth.Logout)
	auth.GET("/me", handlers.Auth.Me)

	projects := api.Group("/projects")
	projects.POST("", handlers.Project.Create)
	projects.GET("", handlers.Project.List)
	projects.GET("/:project_id", handlers.Project.Detail)
	projects.GET("/:project_id/events", handlers.ProjectStream.Stream)

	return router
}
