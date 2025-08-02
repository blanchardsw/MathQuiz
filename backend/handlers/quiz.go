package handlers

import (
	"encoding/json"
	"log"
	"mental-math-trainer/backend/utils"
	"net/http"
)

// HandleQuiz serves a new question based on difficulty
func HandleQuiz(w http.ResponseWriter, r *http.Request) {
	// Get difficulty from query string, e.g. /api/quiz?difficulty=hard
	difficulty := r.URL.Query().Get("difficulty")
	if difficulty == "" {
		difficulty = "normal" // default
	}

	question := utils.GenerateQuestion(difficulty)

	// Get session data
	sessionData, sessionID, err := getSession(r)
	if err != nil {
		http.Error(w, "Session not found", http.StatusUnauthorized)
		return
	}

	// Set session cookie if it's a new session
	if _, err := r.Cookie("session_id"); err != nil {
		http.SetCookie(w, &http.Cookie{
			Name:     "session_id",
			Value:    sessionID,
			Path:     "/",
			HttpOnly: false,
			SameSite: http.SameSiteLaxMode,
			Secure:   false,
		})
	}

	// Store the current question, answer, and difficulty in session for answer validation
	sessionData.CurrentQuestion = question
	sessionData.CurrentDifficulty = difficulty

	log.Printf("Generated [%s] question for session %s: %d %s %d",
		difficulty, sessionID[:8], question.Operand1, question.Operator, question.Operand2)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(question)
}
