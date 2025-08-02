package handlers

import (
	"log"
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
		// Create new session
		sessionID := generateSessionID()
		sessionData := &SessionData{
			HighScores: make(map[string]int),
		}
		sessionsMutex.Lock()
		sessions[sessionID] = sessionData
		sessionsMutex.Unlock()

		http.SetCookie(w, &http.Cookie{
			Name:     "session_id",
			Value:    sessionID,
			Path:     "/",
			HttpOnly: false,
			SameSite: http.SameSiteLaxMode,
			Secure:   false,
		})

		// ✅ Use sessionID in a log to silence staticcheck
		log.Printf("Created new session: %s", sessionID)
	}

	w.WriteHeader(http.StatusOK)
}
