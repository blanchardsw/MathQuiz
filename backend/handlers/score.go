package handlers

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"log"
	"net/http"
	"sync"
	"mental-math-trainer/backend/models"
)

// SessionData holds all session-specific data
type SessionData struct {
	Score             int                `json:"score"`
	CurrentQuestion   models.Question    `json:"currentQuestion"`
	CurrentAnswer     int                `json:"currentAnswer"`     // Store answer separately since Question.Answer has json:"-"
	CurrentDifficulty string             `json:"currentDifficulty"`
	HighScores        map[string]int     `json:"highScores"`
}

// Global session storage with mutex for thread safety
var (
	sessions = make(map[string]*SessionData)
	sessionsMutex = sync.RWMutex{}
)

// generateSessionID creates a new random session ID
func generateSessionID() string {
	bytes := make([]byte, 16)
	rand.Read(bytes)
	return hex.EncodeToString(bytes)
}

// getOrCreateSession retrieves or creates a session for the request
func getOrCreateSession(r *http.Request) (string, *SessionData) {
	sessionsMutex.Lock()
	defer sessionsMutex.Unlock()

	// Try to get session ID from cookie
	cookie, err := r.Cookie("session_id")
	var sessionID string
	if err != nil || cookie.Value == "" {
		// Create new session
		sessionID = generateSessionID()
	} else {
		sessionID = cookie.Value
	}

	// Get or create session data
	sessionData, exists := sessions[sessionID]
	if !exists {
		sessionData = &SessionData{
			Score:             0,
			CurrentQuestion:   models.Question{},
			CurrentAnswer:     0,
			CurrentDifficulty: "",
			HighScores: map[string]int{
				"easy":   0,
				"normal": 0,
				"hard":   0,
			},
		}
		sessions[sessionID] = sessionData
	}

	return sessionID, sessionData
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

	// Check if current score is a new record (but not if previous high score was 0)
	isNewRecord := false
	if sessionData.CurrentDifficulty != "" && sessionData.Score > sessionData.HighScores[sessionData.CurrentDifficulty] && sessionData.HighScores[sessionData.CurrentDifficulty] > 0 {
		sessionData.HighScores[sessionData.CurrentDifficulty] = sessionData.Score
		isNewRecord = true
	} else if sessionData.CurrentDifficulty != "" && sessionData.Score > sessionData.HighScores[sessionData.CurrentDifficulty] {
		// Update high score but don't celebrate if previous was 0
		sessionData.HighScores[sessionData.CurrentDifficulty] = sessionData.Score
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
