package types

import "time"

type ResponsePredict struct {
	ID        uint      `json:"id"`
	UserID    uint      `json:"user_id"`
	Predict   string    `json:"predict"`
	CreatedAt time.Time `json:"created_at"`
}
