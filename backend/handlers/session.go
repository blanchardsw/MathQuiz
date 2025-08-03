package handlers

import (
	"encoding/json"
	"log"
	"mental-math-trainer/backend/utils"
	"net/http"
)

func HandleInitSession(w http.ResponseWriter, r *http.Request) {
	log.Println("Handling /api/init-session")
	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	// Check if session already exists
	_, _, err := getSession(r)
	if err != nil {
		log.Printf("No valid session found: %v. Creating new session.", err)
		sessionID := generateSessionID()
		sessionData := &SessionData{
			HighScores: make(map[string]int),
		}
		sessionsMutex.Lock()
		sessions[sessionID] = sessionData
		sessionsMutex.Unlock()

		// ✅ Create JWT with sessionID
		token, err := utils.GenerateJWT(sessionID)
		if err != nil {
			http.Error(w, "Failed to create token", http.StatusInternalServerError)
			return
		}

		// ✅ Return token in response body
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"token": token})

		log.Printf("Created new session: %s", sessionID)
		return
	}

	w.WriteHeader(http.StatusOK)
}
