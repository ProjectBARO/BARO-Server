package types

import "time"

type ResponseReportSummary struct {
	ID        uint      `json:"id"`
	CreatedAt time.Time `json:"created_at"`
}

type ResponseReport struct {
	ID           uint      `json:"id"`
	UserID       uint      `json:"user_id"`
	AlertCount   int       `json:"alert_count"`
	AnalysisTime int       `json:"analysis_time"`
	Predict      string    `json:"predict"`
	Type         string    `json:"type"`
	CreatedAt    time.Time `json:"created_at"`
}
