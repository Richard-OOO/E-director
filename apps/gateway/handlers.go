package main

import "net/http"

func handleHealth(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, responseEnvelope{
		Code:    0,
		Message: "ok",
		Data: map[string]string{
			"status": "healthy",
		},
	})
}

func handleVersion(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, responseEnvelope{
		Code:    0,
		Message: "ok",
		Data: map[string]string{
			"version": "v1",
		},
	})
}

func handleLogin(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, responseEnvelope{
		Code:    0,
		Message: "login endpoint ready",
		Data: map[string]string{
			"next": "wire auth service",
		},
	})
}

func handleRegister(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, responseEnvelope{
		Code:    0,
		Message: "register endpoint ready",
		Data: map[string]string{
			"next": "wire auth service",
		},
	})
}

func handleMe(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, responseEnvelope{
		Code:    0,
		Message: "me endpoint ready",
		Data: map[string]string{
			"next": "wire auth service",
		},
	})
}
