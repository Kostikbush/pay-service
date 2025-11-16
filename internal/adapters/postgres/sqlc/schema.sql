CREATE TYPE subscription_status AS ENUM ('active','disable','cancel','waiting');

CREATE TABLE subscriptions (
  id              TEXT PRIMARY KEY,
  user_id         TEXT NOT NULL UNIQUE,
  status          subscription_status NOT NULL,
  last_pay_date   TIMESTAMPTZ,
  created_at      TIMESTAMPTZ NOT NULL DEFAULT now(),
  payments_count  INTEGER NOT NULL DEFAULT 0,
  access_until    TIMESTAMPTZ,
  rebill_id       TEXT,
  pm_active       BOOLEAN NOT NULL DEFAULT false,
  card_brand      TEXT,
  last4           TEXT
);
