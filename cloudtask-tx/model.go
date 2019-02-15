package tasktx

import "time"

type Sample struct {
	ID        string    `json:"id"`
	Value     float64   `json:"value"`
	CreatedAt time.Time `json:"createdAt"`
}

type TxStatus struct {
	ID        string    `json:"id"`
	CreatedAt time.Time `json:"createdAt"`
}
