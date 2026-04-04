
CREATE TABLE IF NOT EXISTS roles (
    id SERIAL PRIMARY KEY,
    name VARCHAR(50) UNIQUE NOT NULL,
    description TEXT
);


INSERT INTO roles (name, description) VALUES 
('admin', 'System Administrator'),
('doctor', 'Medical Doctor'),
('patient', 'Hospital Patient')
ON CONFLICT (name) DO NOTHING;


ALTER TABLE users ADD COLUMN IF NOT EXISTS role_id INTEGER REFERENCES roles(id);


CREATE TABLE IF NOT EXISTS refresh_tokens (
    id SERIAL PRIMARY KEY,
    user_id VARCHAR(50) REFERENCES users(id),
    refresh_token TEXT NOT NULL,
    expires_at TIMESTAMP NOT NULL,
    revoked_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT NOW()
);


CREATE TABLE IF NOT EXISTS medical_services (
    id SERIAL PRIMARY KEY,
    service_code VARCHAR(50) UNIQUE NOT NULL,
    service_name VARCHAR(255) NOT NULL,
    base_price DECIMAL(10, 2) NOT NULL DEFAULT 0.00,
    description TEXT,
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS appointments (
    id VARCHAR(255) PRIMARY KEY,
    patient_id VARCHAR(255),
    doctor_id VARCHAR(255),
    service_id VARCHAR(255),
    start_time TIMESTAMP,
    status VARCHAR(50),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

ALTER TABLE appointments 
ADD COLUMN IF NOT EXISTS base_price_at_booking NUMERIC(15,2),
ADD COLUMN IF NOT EXISTS surcharge_at_booking NUMERIC(15,2),
ADD COLUMN IF NOT EXISTS total_price_at_booking NUMERIC(15,2),
ADD COLUMN IF NOT EXISTS applied_policy_snapshot TEXT;


CREATE TABLE IF NOT EXISTS cancellation_policies (
    id SERIAL PRIMARY KEY,
    policy_name VARCHAR(150) UNIQUE NOT NULL,
    hours_before INTEGER NOT NULL,           
    refund_percentage NUMERIC(5,2) NOT NULL,  
    description TEXT,
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

INSERT INTO cancellation_policies (policy_name, hours_before, refund_percentage, description) VALUES
('Free Cancellation (24h+)', 24, 100.00, '100% refund for cancellations made at least 24 hours in advance'),
('Late Cancellation (12h-24h)', 12, 50.00, '50% refund for cancellations made between 12 and 24 hours before the appointment'),
('Last Minute (Under 12h)', 0, 0.00, 'No refund for cancellations made less than 12 hours before the appointment')
ON CONFLICT (policy_name) DO NOTHING;

ALTER TABLE cancellation_policies
ADD COLUMN IF NOT EXISTS version INTEGER NOT NULL DEFAULT 1,
ADD COLUMN IF NOT EXISTS effective_from TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
ADD COLUMN IF NOT EXISTS effective_to TIMESTAMP;

CREATE INDEX IF NOT EXISTS idx_policy_effective_dates ON cancellation_policies(effective_from, effective_to);