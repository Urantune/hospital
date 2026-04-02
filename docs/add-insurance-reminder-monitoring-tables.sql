CREATE EXTENSION IF NOT EXISTS "pgcrypto";

ALTER TABLE users
ADD COLUMN IF NOT EXISTS device_token VARCHAR(500);

CREATE TABLE IF NOT EXISTS insurance_plans (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL,
    provider_name VARCHAR(255),
    coverage_percentage NUMERIC(5,2) NOT NULL CHECK (coverage_percentage >= 0 AND coverage_percentage <= 100),
    status VARCHAR(20) NOT NULL DEFAULT 'active' CHECK (status IN ('active', 'inactive')),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS insurance_service_coverage (
    insurance_plan_id UUID NOT NULL REFERENCES insurance_plans(id) ON DELETE CASCADE,
    service_id UUID NOT NULL REFERENCES services(id) ON DELETE CASCADE,
    custom_coverage_percentage NUMERIC(5,2) NOT NULL CHECK (custom_coverage_percentage >= 0 AND custom_coverage_percentage <= 100),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (insurance_plan_id, service_id)
);

CREATE TABLE IF NOT EXISTS appointment_reminders (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    appointment_id UUID NOT NULL REFERENCES appointments(id) ON DELETE CASCADE,
    reminder_type VARCHAR(20) NOT NULL CHECK (reminder_type IN ('EMAIL', 'SMS', 'PUSH', 'email', 'sms', 'push')),
    scheduled_at TIMESTAMP NOT NULL,
    sent_at TIMESTAMP,
    status VARCHAR(20) NOT NULL DEFAULT 'pending' CHECK (status IN ('pending', 'sent', 'failed', 'cancelled')),
    retry_count INT NOT NULL DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_reminders_appointment ON appointment_reminders(appointment_id);
CREATE INDEX IF NOT EXISTS idx_reminders_status_schedule ON appointment_reminders(status, scheduled_at);

CREATE TABLE IF NOT EXISTS system_failures (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    resource_type VARCHAR(100) NOT NULL,
    resource_id VARCHAR(255),
    failure_reason VARCHAR(255) NOT NULL,
    error_message TEXT,
    details JSONB NOT NULL DEFAULT '{}'::jsonb,
    severity VARCHAR(20) NOT NULL DEFAULT 'ERROR',
    resolved BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_system_failures_created_at ON system_failures(created_at);
CREATE INDEX IF NOT EXISTS idx_system_failures_resource_type ON system_failures(resource_type);
CREATE INDEX IF NOT EXISTS idx_system_failures_resolved ON system_failures(resolved);

CREATE TABLE IF NOT EXISTS system_metrics (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    metric_name VARCHAR(150) NOT NULL,
    metric_value NUMERIC(18,4) NOT NULL,
    dimension_type VARCHAR(100),
    dimension_value VARCHAR(255),
    dimension_date DATE,
    dimension_hour INT,
    recorded_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_system_metrics_name_recorded_at ON system_metrics(metric_name, recorded_at);

CREATE TABLE IF NOT EXISTS api_performance_metrics (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    endpoint VARCHAR(255) NOT NULL,
    method VARCHAR(10) NOT NULL,
    response_time_ms INT NOT NULL,
    status_code INT NOT NULL,
    success BOOLEAN NOT NULL,
    error_message TEXT,
    recorded_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_api_performance_endpoint_recorded_at ON api_performance_metrics(endpoint, recorded_at);
