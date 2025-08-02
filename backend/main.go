package main

import (
	"log"
	"mental-math-trainer/backend/handlers"
	"net/http"
)

func corsMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")
		if origin == "http://localhost:3000" || origin == "https://gomathquiz.netlify.app" {
			w.Header().Set("Access-Control-Allow-Origin", origin)
		}

		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		w.Header().Set("Access-Control-Allow-Credentials", "true")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next(w, r)
	}
}

func main() {
	http.HandleFunc("/api/quiz", corsMiddleware(handlers.HandleQuiz))
	http.HandleFunc("/api/question", corsMiddleware(handlers.HandleGenerateQuestion))
	http.HandleFunc("/api/answer", corsMiddleware(handlers.HandleAnswer))
	http.HandleFunc("/api/score", corsMiddleware(handlers.HandleScore))
	http.HandleFunc("/api/reset-score", corsMiddleware(handlers.HandleResetScore))
	http.HandleFunc("/api/init-session", corsMiddleware(handlers.HandleInitSession))

	log.Println("Server running on port 4000")
	err := http.ListenAndServe(":4000", nil)
	if err != nil {
		log.Fatal("Server failed:", err)
	}
}
