package main

import (
	"log"
	"net/http"
	"mental-math-trainer/backend/handlers"
)

// corsMiddleware adds CORS headers to allow frontend requests
func corsMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "https://gomathquiz.netlify.app")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		
		// Handle preflight requests
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		
		next(w, r)
	}
}

func main() {
	http.HandleFunc("/api/quiz", corsMiddleware(handlers.HandleQuiz))
	http.HandleFunc("/api/answer", corsMiddleware(handlers.HandleAnswer))
	http.HandleFunc("/api/score", corsMiddleware(handlers.HandleScore))
	http.HandleFunc("/api/reset-score", corsMiddleware(handlers.HandleResetScore))

	port := "4000"
	log.Println("Server running on port " + port)
	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		log.Fatal("Server failed:", err)
	}
}