package types

type ResponseVideo struct {
	VideoID      string `json:"video_id"`
	Title        string `json:"title"`
	ThumbnailUrl string `json:"thumbnail_url"`
	Category     string `json:"category"`
}
