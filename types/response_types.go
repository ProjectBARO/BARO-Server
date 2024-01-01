package types

type Response struct {
	Status  int         `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

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

type ResponseVideo struct {
	VideoID      string `json:"video_id"`
	Title        string `json:"title"`
	ThumbnailUrl string `json:"thumbnail_url"`
}
