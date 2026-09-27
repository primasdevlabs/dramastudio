package domain

import "time"

type FactVersion struct {
	FactID    string    `json:"fact_id"`
	Version   int       `json:"version"`
	UpdatedAt time.Time `json:"updated_at"`
}
