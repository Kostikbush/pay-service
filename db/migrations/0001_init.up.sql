CREATE TYPE subscription_status AS ENUM ('active','disable','cancel','waiting');

CREATE TABLE IF NOT EXISTS subscriptions (
  id              TEXT PRIMARY KEY,
  user_id         TEXT NOT NULL UNIQUE,               -- одна подписка на пользователя
  status          subscription_status NOT NULL,
  last_pay_date   TIMESTAMPTZ,
  created_at      TIMESTAMPTZ NOT NULL DEFAULT now(),
  payments_count  INTEGER NOT NULL DEFAULT 0,
  access_until    TIMESTAMPTZ,

  rebill_id       TEXT,                               -- snapshot платёжного метода
  pm_active       BOOLEAN NOT NULL DEFAULT false,
  card_brand      TEXT,
  last4           TEXT
);

CREATE INDEX IF NOT EXISTS idx_subscriptions_status ON subscriptions(status);