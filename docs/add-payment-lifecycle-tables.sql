CREATE TABLE IF NOT EXISTS payment_callback_receipts (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    callback_id VARCHAR(255) UNIQUE NOT NULL,
    payment_id UUID REFERENCES payments(id) ON DELETE CASCADE,
    gateway VARCHAR(50),
    transaction_code VARCHAR(255),
    status VARCHAR(30) NOT NULL,
    raw_payload JSONB,
    notes TEXT,
    received_at TIMESTAMP DEFAULT now(),
    processed_at TIMESTAMP DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_payment_callback_payment ON payment_callback_receipts(payment_id);
CREATE INDEX IF NOT EXISTS idx_payment_callback_status ON payment_callback_receipts(status);

CREATE TABLE IF NOT EXISTS domain_events (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    event_type VARCHAR(150) NOT NULL,
    aggregate_type VARCHAR(100) NOT NULL,
    aggregate_id UUID NOT NULL,
    payload JSONB NOT NULL,
    status VARCHAR(30) NOT NULL DEFAULT 'pending',
    created_at TIMESTAMP DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_domain_event_aggregate ON domain_events(aggregate_type, aggregate_id);
CREATE INDEX IF NOT EXISTS idx_domain_event_status ON domain_events(status);
