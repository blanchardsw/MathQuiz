package main

import (
	"log"
	"mental-math-trainer/handlers"
	"net/http"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

func getAllowedOrigins() map[string]bool {
	raw := os.Getenv("ALLOWED_ORIGINS")
	allowed := make(map[string]bool)
	for _, origin := range strings.Split(raw, ",") {
		origin = strings.TrimSpace(strings.TrimRight(origin, "/"))
		allowed[origin] = true
	}
	return allowed
}

func corsMiddleware(next http.Handler) http.Handler {
	allowedOrigins := getAllowedOrigins()
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := strings.TrimSpace(strings.TrimRight(r.Header.Get("Origin"), "/"))
		log.Println("Incoming Origin:", origin)
		log.Println("Allowed Origins Map:", allowedOrigins)

		if allowedOrigins[origin] {
			log.Println("Origin matched. Setting CORS headers.")
			w.Header().Set("Access-Control-Allow-Origin", origin)
			w.Header().Set("Access-Control-Allow-Credentials", "true")
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		} else {
			log.Println("Origin not allowed:", origin)
		}

		if r.Method == "OPTIONS" {
			log.Println("Handling preflight OPTIONS request for:", r.URL.Path)
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func main() {
	if !strings.Contains(os.Getenv("GO_ENV"), "prod") {
		err := godotenv.Load()
		log.Println("Loaded ALLOWED_ORIGINS:", os.Getenv("ALLOWED_ORIGINS"))
		if err != nil {
			log.Fatal("Error loading .env file")
		}
	}
	port := os.Getenv("PORT")
	if port == "" {
		port = "4000" // fallback for local dev
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "OPTIONS" {
			log.Println("Global OPTIONS handler triggered for:", r.URL.Path)
			// Let middleware handle CORS headers
			w.WriteHeader(http.StatusOK)
			return
		}
		http.NotFound(w, r)
	})

	http.Handle("/api/quiz", corsMiddleware(http.HandlerFunc(handlers.HandleQuiz)))
	http.Handle("/api/question", corsMiddleware(http.HandlerFunc(handlers.HandleGenerateQuestion)))
	http.Handle("/api/answer", corsMiddleware(http.HandlerFunc(handlers.HandleAnswer)))
	http.Handle("/api/score", corsMiddleware(http.HandlerFunc(handlers.HandleScore)))
	http.Handle("/api/reset-score", corsMiddleware(http.HandlerFunc(handlers.HandleResetScore)))
	http.Handle("/api/init-session", corsMiddleware(http.HandlerFunc(handlers.HandleInitSession)))
	http.Handle("/api/health", corsMiddleware(http.HandlerFunc(handlers.HandleHealth)))

	log.Println("Server running on port", port)
	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		log.Fatal("Server failed:", err)
	}
}
