package types

import "github.com/go-playground/validator/v10"

var validate *validator.Validate

func init() {
	validate = validator.New()
}

type RequestPredict struct {
	VideoURL string `json:"video_url" validate:"required"`
}

func (r *RequestPredict) Validate() error {
	return validate.Struct(r)
}
