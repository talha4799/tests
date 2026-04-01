package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

// ─── Response helpers ────────────────────────────────────────────────────────

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}

func methodNotAllowed(w http.ResponseWriter, allowed string) {
	w.Header().Set("Allow", allowed)
	writeJSON(w, http.StatusMethodNotAllowed, map[string]string{
		"error":   "Method Not Allowed",
		"allowed": allowed,
	})
}

// ─── Middleware ───────────────────────────────────────────────────────────────

func corsMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		next(w, r)
	}
}

func onlyAllow(method string, next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != method {
			methodNotAllowed(w, method)
			return
		}
		next(w, r)
	}
}

// ─── Handlers ────────────────────────────────────────────────────────────────

// GET /health  — is the server alive?
func healthHandler(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]any{
		"status":    "ok",
		"message":   "Server is healthy",
		"timestamp": time.Now().UTC().Format(time.RFC3339),
	})
}

// GET /check  — detailed server check
func checkHandler(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]any{
		"status":    "running",
		"uptime":    fmt.Sprintf("%s", time.Since(startTime).Round(time.Second)),
		"server":    "Go HTTP Server",
		"version":   "1.0.0",
		"checkedAt": time.Now().UTC().Format(time.RFC3339),
	})
}

// GET /weather  — fake weather data (GET only)
func weatherHandler(w http.ResponseWriter, r *http.Request) {
	city := r.URL.Query().Get("city")
	if city == "" {
		city = "Faisalabad"
	}

	// Fake but realistic weather payload
	writeJSON(w, http.StatusOK, map[string]any{
		"city":    city,
		"country": "PK",
		"temperature": map[string]any{
			"celsius":    32,
			"fahrenheit": 89.6,
			"feels_like": 35,
		},
		"condition": "Partly Cloudy",
		"humidity":  58,
		"wind": map[string]any{
			"speed":     14,
			"direction": "NW",
			"unit":      "km/h",
		},
		"visibility": "10 km",
		"uv_index":   7,
		"sunrise":    "06:12 AM",
		"sunset":     "06:48 PM",
		"fetchedAt":  time.Now().UTC().Format(time.RFC3339),
	})
}

// ─── 404 ─────────────────────────────────────────────────────────────────────

func notFoundHandler(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusNotFound, map[string]string{
		"error": "Route not found",
		"path":  r.URL.Path,
	})
}

// ─── Bootstrap ───────────────────────────────────────────────────────────────

var startTime = time.Now()

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("/health", corsMiddleware(onlyAllow(http.MethodGet, healthHandler)))
	mux.HandleFunc("/check", corsMiddleware(onlyAllow(http.MethodGet, checkHandler)))
	mux.HandleFunc("/weather", corsMiddleware(onlyAllow(http.MethodGet, weatherHandler)))
	mux.HandleFunc("/", corsMiddleware(notFoundHandler))

	port := ":8088"
	fmt.Printf("🚀 Server running on http://localhost%s\n", port)
	fmt.Println("   GET /health   → server health status")
	fmt.Println("   GET /check    → detailed server check")
	fmt.Println("   GET /weather  → fake weather JSON (add ?city=Lahore)")

	log.Fatal(http.ListenAndServe(port, mux))
}
