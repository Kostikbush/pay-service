-- name: UpsertSubscription :exec
INSERT INTO subscriptions (
  id, user_id, status, last_pay_date, created_at, payments_count, access_until,
  rebill_id, pm_active, card_brand, last4
) VALUES (
  sqlc.arg(id),
  sqlc.arg(user_id),
  sqlc.arg(status)::subscription_status,
  sqlc.arg(last_pay_date),
  sqlc.arg(created_at),
  sqlc.arg(payments_count),
  sqlc.arg(access_until),
  sqlc.arg(rebill_id),
  sqlc.arg(pm_active),
  sqlc.arg(card_brand),
  sqlc.arg(last4)
)
ON CONFLICT (user_id) DO UPDATE SET
  id=EXCLUDED.id,
  status=EXCLUDED.status,
  last_pay_date=EXCLUDED.last_pay_date,
  payments_count=EXCLUDED.payments_count,
  access_until=EXCLUDED.access_until,
  rebill_id=EXCLUDED.rebill_id,
  pm_active=EXCLUDED.pm_active,
  card_brand=EXCLUDED.card_brand,
  last4=EXCLUDED.last4;

-- name: GetSubscriptionByUserID :one
SELECT
  id, user_id, status, last_pay_date, created_at, payments_count, access_until,
  rebill_id, pm_active, card_brand, last4
FROM subscriptions
WHERE user_id = sqlc.arg(user_id);

-- name: SetRebillID :exec
UPDATE subscriptions
SET rebill_id = sqlc.arg(rebill_id),
    pm_active = true,
    card_brand = sqlc.arg(card_brand),
    last4 = sqlc.arg(last4)
WHERE user_id = sqlc.arg(user_id);