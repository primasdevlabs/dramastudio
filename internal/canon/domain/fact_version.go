package domain

import (
	"encoding/json"
	"time"
)

type FactVersion struct {
	FactID    string          `json:"fact_id"`
	Version   int             `json:"version"`
	Value     json.RawMessage `json:"value"`
	UpdatedAt time.Time       `json:"updated_at"`
}
