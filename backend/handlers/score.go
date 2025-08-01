package handlers

import (
	"encoding/json"
	"net/http"
	"mental-math-trainer/backend/models"
)

var score int
var currentQuestion models.Question // Store the current question server-side
var currentDifficulty string

// High scores per difficulty (in-memory storage)
var highScores = map[string]int{
	"easy":   0,
	"normal": 0,
	"hard":   0,
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

	// Compare against the server-side stored answer
	correct := req.UserAnswer == currentQuestion.Answer
	if correct {
		score++
	}

	resp := models.AnswerResponse{
		Correct:       correct,
		CorrectAnswer: currentQuestion.Answer,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

// HandleScore returns the current score and high scores
func HandleScore(w http.ResponseWriter, r *http.Request) {
	// Check if current score is a new record (but not if previous high score was 0)
	isNewRecord := false
	if currentDifficulty != "" && score > highScores[currentDifficulty] && highScores[currentDifficulty] > 0 {
		highScores[currentDifficulty] = score
		isNewRecord = true
	} else if currentDifficulty != "" && score > highScores[currentDifficulty] {
		// Update high score but don't celebrate if previous was 0
		highScores[currentDifficulty] = score
	}

	resp := models.ScoreResponse{
		CurrentScore: score,
		HighScores:   highScores,
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

	// Update high score if current score is higher
	if score > highScores[req.Difficulty] {
		highScores[req.Difficulty] = score
	}

	// Reset current score
	score = 0
	currentDifficulty = ""

	resp := models.ScoreResponse{
		CurrentScore: score,
		HighScores:   highScores,
		IsNewRecord:  false,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}
