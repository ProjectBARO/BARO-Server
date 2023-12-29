package types

type RequestCreateUser struct {
	Name   string `json:"name" binding:"required"`
	Email  string `json:"email" binding:"required"`
	Age    int    `json:"age"`
	Gender string `json:"gender"`
}