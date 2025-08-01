package models

// AnswerRequest is the payload from the frontend
type AnswerRequest struct {
	UserAnswer int      `json:"userAnswer"`
	Question   Question `json:"question"`
}

// AnswerResponse is the feedback sent to the frontend
type AnswerResponse struct {
	Correct       bool `json:"correct"`
	CorrectAnswer int  `json:"correctAnswer"`
}
