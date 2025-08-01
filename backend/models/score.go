package models

// ScoreResponse represents the current session score and high scores
type ScoreResponse struct {
	CurrentScore int                    `json:"currentScore"`
	HighScores   map[string]int         `json:"highScores"`
	IsNewRecord  bool                   `json:"isNewRecord"`
}

// ResetRequest represents a request to reset the current score
type ResetRequest struct {
	Difficulty string `json:"difficulty"`
}
