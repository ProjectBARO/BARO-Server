package types

import "github.com/go-playground/validator/v10"

var validate *validator.Validate

func init() {
	validate = validator.New()
}

type RequestCreateUser struct {
	Name   string `json:"name" validate:"required"`
	Email  string `json:"email" validate:"required,email"`
	Age    int    `json:"age"`
	Gender string `json:"gender"`
}

func (r *RequestCreateUser) Validate() error {
	return validate.Struct(r)
}

type RequestUpdateUser struct {
	Nickname string `json:"nickname"`
	Age      int    `json:"age"`
	Gender   string `json:"gender"`
}

func (r *RequestUpdateUser) Validate() error {
	return validate.Struct(r)
}
