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
	
	// Store the current question and difficulty server-side for answer validation
	currentQuestion = question
	currentDifficulty = difficulty

	log.Printf("Generated [%s] question: %d %s %d = %d",
		difficulty, question.Operand1, question.Operator, question.Operand2, question.Answer)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(question)
}
