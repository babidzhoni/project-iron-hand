package main

type ML1ScoredMessage struct {
	CaseID string `json:"case_id"`
	Scores struct {
		Score  float64 `json:"score"`
		Memory int     `json:"memory"`
		Logs   string  `json:"logs"`
	} `json:"scores"`
}
