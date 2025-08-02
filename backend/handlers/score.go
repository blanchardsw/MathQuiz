package handlers

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
	"log"
	"mental-math-trainer/backend/models"
	"net/http"
	"sync"
)

// SessionData holds all session-specific data
type SessionData struct {
	Score             int             `json:"score"`
	CurrentQuestion   models.Question `json:"currentQuestion"`
	CurrentAnswer     int             `json:"currentAnswer"` // Store answer separately since Question.Answer has json:"-"
	CurrentDifficulty string          `json:"currentDifficulty"`
	HighScores        map[string]int  `json:"highScores"`
}

// Global session storage with mutex for thread safety
var (
	sessions      = make(map[string]*SessionData)
	sessionsMutex = sync.RWMutex{}
)

// generateSessionID creates a new random session ID
func generateSessionID() string {
	bytes := make([]byte, 16)
	rand.Read(bytes)
	return hex.EncodeToString(bytes)
}

func getSession(r *http.Request) (*SessionData, string, error) {
	cookie, err := r.Cookie("session_id")
	if err != nil {
		return nil, "", err
	}

	sessionsMutex.RLock()
	defer sessionsMutex.RUnlock()

	sessionData, exists := sessions[cookie.Value]
	if !exists {
		return nil, cookie.Value, errors.New("session not found")
	}

	return sessionData, cookie.Value, nil
}

// HandleAnswer checks the user's answer and updates score
func HandleAnswer(w http.ResponseWriter, r *http.Request) {
	var req struct {
		UserAnswer int `json:"userAnswer"`
	}

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

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

	log.Printf("DEBUG: HandleAnswer - Session ID: %s, CurrentAnswer: %d, UserAnswer: %d", sessionID[:8], sessionData.CurrentAnswer, req.UserAnswer)

	// Compare against the session-stored answer
	correct := req.UserAnswer == sessionData.CurrentAnswer
	if correct {
		sessionData.Score++
	}

	resp := models.AnswerResponse{
		Correct:       correct,
		CorrectAnswer: sessionData.CurrentAnswer,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

// HandleScore returns the current score and high scores
func HandleScore(w http.ResponseWriter, r *http.Request) {
	sessionData, sessionID, err := getSession(r)
	log.Println("Handling /api/score")
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

	// Check if current score is a new record (but not if previous high score was 0)
	isNewRecord := false
	// Ensure high score is initialized
	if sessionData.CurrentDifficulty != "" {
		if _, exists := sessionData.HighScores[sessionData.CurrentDifficulty]; !exists {
			sessionData.HighScores[sessionData.CurrentDifficulty] = 0
		}

		// Check if current score is a new record
		if sessionData.Score > sessionData.HighScores[sessionData.CurrentDifficulty] {
			isNewRecord = sessionData.HighScores[sessionData.CurrentDifficulty] > 0
			sessionData.HighScores[sessionData.CurrentDifficulty] = sessionData.Score
		}
	}

	resp := models.ScoreResponse{
		CurrentScore: sessionData.Score,
		HighScores:   sessionData.HighScores,
		IsNewRecord:  isNewRecord,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

// HandleResetScore resets the current score when quiz ends
func HandleResetScore(w http.ResponseWriter, r *http.Request) {
	var req models.ResetRequest

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

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

	// Update high score if current score is higher
	if sessionData.Score > sessionData.HighScores[req.Difficulty] {
		sessionData.HighScores[req.Difficulty] = sessionData.Score
	}

	// Reset current score
	sessionData.Score = 0
	sessionData.CurrentDifficulty = ""

	resp := models.ScoreResponse{
		CurrentScore: sessionData.Score,
		HighScores:   sessionData.HighScores,
		IsNewRecord:  false,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}
