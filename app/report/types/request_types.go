package types

import "github.com/go-playground/validator/v10"

var validate *validator.Validate

func init() {
	validate = validator.New()
}

type RequestAnalysis struct {
	VideoURL     string `json:"video_url" validate:"required"`
	AlertCount   uint   `json:"alert_count" validate:"required"`
	AnalysisTime uint   `json:"analysis_time" validate:"required"`
	Type         string `json:"type" validate:"required"`
}

type RequestReportSummary struct {
	YearAndMonth string `json:"year_and_month" validate:"required"`
}

func (r *RequestAnalysis) Validate() error {
	return validate.Struct(r)
}
