package handlers

import (
	"encoding/json"
	"math/rand"
	"mental-math-trainer/models"
	"net/http"
	"strings"
	"time"
)

var rng = rand.New(rand.NewSource(time.Now().UnixNano()))

// generateQuestion creates a new question based on difficulty
func generateQuestion(difficulty string) models.Question {
	var min, max int
	var operators []string

	switch difficulty {
	case "easy":
		min, max = 1, 10
		operators = []string{"+", "-"}
	case "normal":
		min, max = 10, 50
		operators = []string{"+", "-", "*"}
	case "hard":
		min, max = 50, 100
		operators = []string{"+", "-", "*", "/"}
	default:
		min, max = 1, 10
		operators = []string{"+"}
	}

	op1 := rng.Intn(max-min+1) + min
	op2 := rng.Intn(max-min+1) + min
	operator := operators[rng.Intn(len(operators))]

	// Avoid divide-by-zero
	if operator == "/" && op2 == 0 {
		op2 = 1
	}

	return models.Question{
		Operand1:   op1,
		Operand2:   op2,
		Operator:   operator,
		Difficulty: difficulty,
	}
}

// HandleGenerateQuestion creates a new question and stores it in session
func HandleGenerateQuestion(w http.ResponseWriter, r *http.Request) {

	// Handle preflight
	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	// Parse difficulty
	difficulty := strings.ToLower(strings.TrimSpace(r.URL.Query().Get("difficulty")))
	if difficulty == "" {
		difficulty = "easy"
	}

	// Generate question
	question := generateQuestion(difficulty)

	// ✅ Compute answer separately
	answer := computeAnswer(question)

	// Retrieve or create session
	sessionData, sessionID, err := getSession(r)
	if err != nil {
		sessionID = generateSessionID()
		sessionData = &SessionData{
			HighScores: make(map[string]int),
		}
		sessionsMutex.Lock()
		sessions[sessionID] = sessionData
		sessionsMutex.Unlock()
	}

	// Store question and answer in session
	sessionData.CurrentQuestion = question
	sessionData.CurrentAnswer = answer
	sessionData.CurrentDifficulty = difficulty

	// Return question without answer
	publicQuestion := struct {
		Operand1   int    `json:"operand1"`
		Operand2   int    `json:"operand2"`
		Operator   string `json:"operator"`
		Difficulty string `json:"difficulty"`
	}{
		Operand1:   question.Operand1,
		Operand2:   question.Operand2,
		Operator:   question.Operator,
		Difficulty: question.Difficulty,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(publicQuestion)
}

func computeAnswer(q models.Question) int {
	switch q.Operator {
	case "+":
		return q.Operand1 + q.Operand2
	case "-":
		return q.Operand1 - q.Operand2
	case "*":
		return q.Operand1 * q.Operand2
	case "/":
		if q.Operand2 != 0 {
			return q.Operand1 / q.Operand2
		}
	}
	return 0 // default fallback
}
