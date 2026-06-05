package main

import (
	"log"
	"net/http"
	"os"
	"time"
)

func main() {
	cfg := loadConfig()
	router := NewRouter(cfg)

	srv := &http.Server{
		Addr:              cfg.ListenAddr,
		Handler:           router,
		ReadHeaderTimeout: 5 * time.Second,
	}

	log.Printf("gateway listening on %s", cfg.ListenAddr)
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("gateway server failed: %v", err)
	}
}

type Config struct {
	ListenAddr     string
	FrontendOrigin string
}

func loadConfig() Config {
	listenAddr := os.Getenv("GATEWAY_ADDR")
	if listenAddr == "" {
		listenAddr = ":8080"
	}

	frontendOrigin := os.Getenv("FRONTEND_ORIGIN")
	if frontendOrigin == "" {
		frontendOrigin = "http://localhost:5173"
	}

	return Config{
		ListenAddr:     listenAddr,
		FrontendOrigin: frontendOrigin,
	}
}
