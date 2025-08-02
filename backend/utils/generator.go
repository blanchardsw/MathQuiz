package utils

import (
	"math/rand"
	"mental-math-trainer/backend/models"
	"time"
)

// GenerateQuestion creates a question based on difficulty
func GenerateQuestion(difficulty string) models.Question {
	rand.Seed(time.Now().UnixNano())

	var min, max int
	var operators []string

	switch difficulty {
	case "easy":
		min, max = 1, 10
		operators = []string{"+"}
	case "normal":
		min, max = 1, 20
		operators = []string{"+", "-"}
	case "hard":
		min, max = 10, 99
		operators = []string{"+", "-", "*"}
	default:
		// fallback to normal
		min, max = 1, 20
		operators = []string{"+", "-"}
	}

	operand1 := rand.Intn(max-min+1) + min
	operand2 := rand.Intn(max-min+1) + min
	operator := operators[rand.Intn(len(operators))]

	return models.Question{
		Operand1:   operand1,
		Operand2:   operand2,
		Operator:   operator,
		Difficulty: difficulty,
	}
}
