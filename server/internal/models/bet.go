package models

type Bet struct {
	ID        int     `json:"id"`
	UserID    int     `json:"user_id"`
	EventID   int     `json:"event_id"`
	Amount    float32 `json:"amount"`
	Outcome   bool    `json:"outcome"`
	Timestamp string  `json:"timestamp"`
	Payout    float32 `json:"payout"`
}
