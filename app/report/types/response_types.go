package types

import "time"

type Landmark struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
}

type LandmarkInfo struct {
	LeftShoulder       Landmark `json:"left_shoulder"`
	LeftEar            Landmark `json:"left_ear"`
	VerticalDistanceCM float64  `json:"vertical_distance_cm"`
	Angle              float64  `json:"angle"`
}

type ResponseReportSummary struct {
	ID        uint      `json:"id"`
	CreatedAt time.Time `json:"created_at"`
}

type ResponseAnalysis struct {
	Result            []int          `json:"result"`
	HunchedRatio      float64        `json:"hunched_ratio"`
	NormalRatio       float64        `json:"normal_ratio"`
	Scores            []float64      `json:"scores"`
	LandmarksInfo     []LandmarkInfo `json:"landmarks_info"`
	StatusFrequencies map[string]int `json:"status_frequencies"`
}

type ResponseReport struct {
	ID                uint      `json:"id"`
	UserID            uint      `json:"user_id"`
	AlertCount        int       `json:"alert_count"`
	AnalysisTime      int       `json:"analysis_time"`
	Type              string    `json:"type"`
	Predict           string    `json:"predict"`
	Score             string    `json:"score"`
	NormalRatio       string    `json:"normal_ratio"`
	NeckAngles        string    `json:"neck_angles"`
	Distances         string    `json:"distances"`
	StatusFrequencies string    `json:"status_frequencies"`
	CreatedAt         time.Time `json:"created_at"`
}

type ResponseRank struct {
	UserID          uint   `json:"user_id"`
	Nickname        string `json:"nickname"`
	Gender          string `json:"gender"`
	Age             int    `json:"age"`
	NormalRatio     string `json:"normal_ratio"`
	AverageScore    string `json:"average_score"`
	AllAverageScore string `json:"all_average_score"`
}
