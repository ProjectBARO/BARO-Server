package types

type ResponseToken struct {
	Token string `json:"token"`
}

type ResponseUser struct {
	ID       uint   `json:"id"`
	Name     string `json:"name"`
	Nickname string `json:"nickname"`
	Email    string `json:"email"`
	Age      int    `json:"age"`
	Gender   string `json:"gender"`
}
