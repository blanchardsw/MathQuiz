package models

// Question represents a single arithmetic problem
type Question struct {
	Operand1   int    `json:"operand1"`
	Operand2   int    `json:"operand2"`
	Operator   string `json:"operator"`
	Difficulty string `json:"difficulty"`
}
