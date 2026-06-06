package main

import (
	"context"
	"log"

	"github.com/Richard-OOO/E-director/apps/gateway/internal/config"
	"github.com/Richard-OOO/E-director/apps/gateway/internal/handler"
	mysqlmodels "github.com/Richard-OOO/E-director/apps/gateway/internal/models/mysql"
	redismodels "github.com/Richard-OOO/E-director/apps/gateway/internal/models/redis"
	"github.com/Richard-OOO/E-director/apps/gateway/internal/server"
	"github.com/Richard-OOO/E-director/apps/gateway/internal/service"
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

	router := server.NewRouter(cfg, server.Handlers{
		System: handler.NewSystemHandler(),
		Auth:   handler.NewAuthHandler(authService, cfg),
	})

	log.Printf("gateway listening on %s", cfg.ListenAddr)
	if err := router.Run(cfg.ListenAddr); err != nil {
		log.Fatalf("gateway server failed: %v", err)
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
