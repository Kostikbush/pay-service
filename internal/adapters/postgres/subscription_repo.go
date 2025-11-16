package postgres

import (
	"context"
	"time"

	gen "pay-service/internal/adapters/postgres/sqlc/gen"
	"pay-service/internal/entities"

	"github.com/jackc/pgx/v5/pgxpool"
)

type SubscriptionRepo struct {
	q *gen.Queries
}

func NewSubscriptionRepo(pool *pgxpool.Pool) *SubscriptionRepo {
	return &SubscriptionRepo{ q: gen.New(pool) }
}

func (r *SubscriptionRepo) Save(ctx context.Context, s *entities.Subscription) error {
	params := gen.UpsertSubscriptionParams{
		ID:            s.ID,
		UserID:        s.UserID,
		Status:        gen.SubscriptionStatus(s.SubscriptionStatus),
		LastPayDate:   nullableTimePtr(s.LastPayDate), // *time.Time
		CreatedAt:     s.CreatedAt,
		PaymentsCount: int32(s.PaymentsCount),
		AccessUntil:   nullableTimePtr(s.AccessUntil), // *time.Time

		RebillID:  nullableStringPtr(s.PaymentMethod.RebillID),
		PmActive:  s.PaymentMethod.Active,
		CardBrand: nullableStringPtr(s.PaymentMethod.Brand),
		Last4:     nullableStringPtr(s.PaymentMethod.Last4),
	}
	return r.q.UpsertSubscription(ctx, params)
}

func (r *SubscriptionRepo) FindByUser(ctx context.Context, userID string) (*entities.Subscription, error) {
	row, err := r.q.GetSubscriptionByUserID(ctx, userID)
	if err != nil {
		return nil, err // при желании заверни в доменную "not found"
	}
	return &entities.Subscription{
		ID:                 row.ID,
		UserID:             row.UserID,
		SubscriptionStatus: entities.SubscriptionStatus(row.Status),
		LastPayDate:        derefTimePtr(&row.LastPayDate.Time),
		CreatedAt:          row.CreatedAt.Time,
		PaymentsCount:      int64(row.PaymentsCount),
		AccessUntil:        derefTimePtr(&row.AccessUntil.Time),
		PaymentMethod: entities.PaymentMethodSnapshot{
			RebillID: derefStringPtr(row.RebillID),
			Active:   row.PmActive,
			Brand:    derefStringPtr(row.CardBrand),
			Last4:    derefStringPtr(row.Last4),
		},
	}, nil
}

func (r *SubscriptionRepo) SetRebillID(ctx context.Context, userID, rebillID, brand, last4 string) error {
	return r.q.SetRebillID(ctx, gen.SetRebillIDParams{
		UserID:   userID,
		RebillID: &rebillID,
		CardBrand: &brand,
		Last4:    &last4,
	})
}

// helpers для nullable
func nullableTimePtr(t time.Time) *time.Time { if t.IsZero() { return nil }; return &t }
func derefTimePtr(t *time.Time) time.Time   { if t == nil { return time.Time{} }; return *t }
func nullableStringPtr(s string) *string    { if s == "" { return nil }; return &s }
func derefStringPtr(s *string) string       { if s == nil { return "" }; return *s }

