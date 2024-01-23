package types

import "github.com/go-playground/validator/v10"

var validate *validator.Validate

func init() {
	validate = validator.New()
}

type RequestAnalysis struct {
	VideoURL     string `json:"video_url" validate:"required"`
	AlertCount   int    `json:"alert_count" validate:"min=0"`
	AnalysisTime int    `json:"analysis_time" validate:"min=0"`
	Type         string `json:"type" validate:"required"`
}

type RequestReportSummary struct {
	YearAndMonth string `json:"year_and_month" validate:"required"`
}

func (r *RequestAnalysis) Validate() error {
	return validate.Struct(r)
}
