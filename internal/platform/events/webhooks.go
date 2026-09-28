package events

import (
	"context"
	"time"

	"github.com/google/uuid"

	"dramastudio/internal/platform/database/postgres"
)

// WebhookDeliveryStore deduplicates provider webhook deliveries against
// platform.webhook_deliveries (§61). Nil-safe: callers skip dedup when the
// store is absent (in-memory dev/test) and rely on terminal-state job
// idempotency instead.
type WebhookDeliveryStore struct {
	q postgres.Querier
}

func NewWebhookDeliveryStore(q postgres.Querier) *WebhookDeliveryStore {
	if q == nil {
		return nil
	}
	return &WebhookDeliveryStore{q: q}
}

// RecordDelivery inserts the delivery; returns false when the id was
// already recorded (replay). A new id is generated when the provider does
// not supply one, so dedup still applies within-process.
func (s *WebhookDeliveryStore) RecordDelivery(ctx context.Context, deliveryID, provider string, payload []byte) (first bool, err error) {
	if deliveryID == "" {
		deliveryID = "dlv_" + uuid.NewString()
	}
	tag, err := s.q.Exec(ctx, `
		INSERT INTO platform.webhook_deliveries (id, provider, received_at, payload)
		VALUES ($1,$2,$3,$4) ON CONFLICT (id) DO NOTHING`,
		deliveryID, provider, time.Now().UTC(), payload)
	if err != nil {
		return false, err
	}
	return tag.RowsAffected() > 0, nil
}
