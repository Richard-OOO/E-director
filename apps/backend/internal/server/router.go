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
	Prompt        *handler.PromptHandler
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
	projects.PATCH("/:project_id/scenes/:scene_id", handlers.Project.UpdateSceneYAML)
	projects.DELETE("/:project_id", handlers.Project.Delete)
	projects.GET("/:project_id/events", handlers.ProjectStream.Stream)
	projects.GET("/:project_id/export/yaml", handlers.Project.ExportYAML)
	projects.GET("/:project_id/export/yaml/combined", handlers.Project.ExportCombinedYAML)

	prompts := api.Group("/prompts")
	prompts.GET("", handlers.Prompt.List)
	prompts.POST("", handlers.Prompt.Create)
	prompts.GET("/:prompt_id", handlers.Prompt.Detail)
	prompts.PATCH("/:prompt_id", handlers.Prompt.Update)
	prompts.DELETE("/:prompt_id", handlers.Prompt.Delete)
	prompts.POST("/:prompt_id/reset", handlers.Prompt.Reset)

	return router
}
