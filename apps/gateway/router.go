package main

import (
	"encoding/json"
	"net/http"
)

type Router struct {
	mux *http.ServeMux
	cfg Config
}

func NewRouter(cfg Config) http.Handler {
	m := http.NewServeMux()
	registerRoutes(m)

	return &Router{
		mux: m,
		cfg: cfg,
	}
}

func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", r.cfg.FrontendOrigin)
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	if req.Method == http.MethodOptions {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	r.mux.ServeHTTP(w, req)
}

func registerRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /health", handleHealth)
	mux.HandleFunc("GET /api/v1/version", handleVersion)
	mux.HandleFunc("POST /api/v1/auth/login", handleLogin)
	mux.HandleFunc("POST /api/v1/auth/register", handleRegister)
	mux.HandleFunc("GET /api/v1/auth/me", handleMe)
}

type responseEnvelope struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func writeJSON(w http.ResponseWriter, status int, payload responseEnvelope) {
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}
