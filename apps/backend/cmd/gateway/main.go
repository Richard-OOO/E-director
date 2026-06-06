package main

import (
	"context"
	"log"

	"github.com/Richard-OOO/E-director/apps/backend/internal/config"
	"github.com/Richard-OOO/E-director/apps/backend/internal/handler"
	mysqlmodels "github.com/Richard-OOO/E-director/apps/backend/internal/models/mysql"
	redismodels "github.com/Richard-OOO/E-director/apps/backend/internal/models/redis"
	"github.com/Richard-OOO/E-director/apps/backend/internal/server"
	"github.com/Richard-OOO/E-director/apps/backend/internal/service"
	"gorm.io/gorm"
)

func main() {
	cfg := config.Load()
	ctx := context.Background()

	db := mustOpenMySQL(cfg)
	rdb := redismodels.NewClient(cfg.RedisAddr, cfg.RedisPassword, cfg.RedisDB)
	if err := redismodels.Ping(ctx, rdb); err != nil {
		log.Fatalf("connect redis failed: %v", err)
	}

	userStore := mysqlmodels.NewUserStore(db)
	generationStore := mysqlmodels.NewGenerationStore(db)
	sessionStore := redismodels.NewSessionStore(rdb, cfg.SessionTTL)
	codeStore := redismodels.NewVerificationCodeStoreWithLimits(rdb, cfg.VerificationTTL, cfg.CodeCooldown, cfg.CodeAttemptTTL)
	mailer := service.NewSMTPMailer(service.SMTPConfig{
		Host:     cfg.SMTPHost,
		Port:     cfg.SMTPPort,
		Username: cfg.SMTPUsername,
		Password: cfg.SMTPPassword,
		From:     cfg.SMTPFrom,
	})
	authService := service.NewAuthService(userStore, sessionStore, codeStore, service.WithMailer(mailer), service.WithMaxCodeAttempts(cfg.CodeMaxAttempts))
	eventBus := service.NewGenerationEventBus()
	generationService := service.NewGenerationService(generationStore, service.NewChapterSplitter(), service.NewChapterPromptBuilder(), service.NewMockChapterGenerator(), eventBus)
	authHandler := handler.NewAuthHandler(authService, cfg)
	projectHandler := handler.NewProjectHandler(authHandler, generationService)
	projectStreamHandler := handler.NewProjectStreamHandler(authHandler, eventBus)

	router := server.NewRouter(cfg, server.Handlers{
		System:        handler.NewSystemHandler(),
		Auth:          authHandler,
		Project:       projectHandler,
		ProjectStream: projectStreamHandler,
	})

	log.Printf("backend listening on %s", cfg.ListenAddr)
	if err := router.Run(cfg.ListenAddr); err != nil {
		log.Fatalf("backend server failed: %v", err)
	}
}

func mustOpenMySQL(cfg config.Config) *gorm.DB {
	if cfg.MySQLDSN == "" {
		log.Fatal("MYSQL_DSN is required")
	}
	db, err := mysqlmodels.Open(cfg.MySQLDSN)
	if err != nil {
		log.Fatalf("connect mysql failed: %v", err)
	}
	if err := mysqlmodels.AutoMigrate(db); err != nil {
		log.Fatalf("migrate mysql failed: %v", err)
	}
	return db
}
