package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"mental-math-trainer/backend/utils"
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
	sessionID, sessionData := getOrCreateSession(r)

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
	sessionData.CurrentAnswer = question.Answer
	sessionData.CurrentDifficulty = difficulty

	log.Printf("Generated [%s] question for session %s: %d %s %d = %d",
		difficulty, sessionID[:8], question.Operand1, question.Operator, question.Operand2, question.Answer)
	log.Printf("DEBUG: Stored answer in session: %d", sessionData.CurrentAnswer)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(question)
}
