-- ==============================
-- MISSING TABLES MIGRATION
-- ==============================

-- ==============================
-- APPOINTMENT STATE HISTORY
-- ==============================

CREATE TABLE appointment_state_history (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    appointment_id UUID,
    from_state VARCHAR(30),
    to_state VARCHAR(30) NOT NULL CHECK (
        to_state IN ('CREATED','PENDING_PAYMENT','CONFIRMED',
                   'IN_PROGRESS','COMPLETED','CANCELLED','NO_SHOW')
    ),
    changed_by UUID,
    reason TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (appointment_id) REFERENCES appointments(id) ON DELETE CASCADE,
    FOREIGN KEY (changed_by) REFERENCES users(id) ON DELETE SET NULL
);

CREATE INDEX idx_appointment_history ON appointment_state_history(appointment_id);
CREATE INDEX idx_history_created_at ON appointment_state_history(created_at);


CREATE TABLE exception_day (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    doctor_id UUID,
    date DATE NOT NULL,
    type VARCHAR(20) NOT NULL CHECK (type IN ('OFF', 'HOLIDAY', 'CLOSURE', 'SPECIAL_HOURS')),
    start_time TIME,
    end_time TIME,
    reason TEXT,
    created_by UUID,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (doctor_id) REFERENCES doctors(id) ON DELETE CASCADE,
    FOREIGN KEY (created_by) REFERENCES users(id) ON DELETE SET NULL
);

CREATE INDEX idx_exception_doctor_date ON exception_day(doctor_id, date);
CREATE INDEX idx_exception_type ON exception_day(type);


CREATE TABLE cms_change_events (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    event_id VARCHAR(255) UNIQUE NOT NULL,
    source VARCHAR(50) NOT NULL,
    entity_type VARCHAR(100) NOT NULL,
    entity_id VARCHAR(255),
    action VARCHAR(50) NOT NULL,
    payload TEXT,
    status VARCHAR(20) NOT NULL CHECK (status IN ('received', 'processing', 'applied', 'failed')),
    error_message TEXT,
    processed_by UUID,
    processed_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (processed_by) REFERENCES users(id) ON DELETE SET NULL
);

CREATE UNIQUE INDEX uq_cms_event_id ON cms_change_events(event_id);
CREATE INDEX idx_cms_status ON cms_change_events(status);
CREATE INDEX idx_cms_entity ON cms_change_events(entity_type, entity_id);


CREATE TABLE audit_logs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID,
    clinic_id UUID,
    action VARCHAR(150) NOT NULL,
    resource VARCHAR(150) NOT NULL,
    resource_id VARCHAR(255),
    description TEXT,
    ip_address VARCHAR(64),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE SET NULL,
    FOREIGN KEY (clinic_id) REFERENCES clinics(id) ON DELETE SET NULL
);

CREATE INDEX idx_audit_user ON audit_logs(user_id);
CREATE INDEX idx_audit_clinic ON audit_logs(clinic_id);
CREATE INDEX idx_audit_action ON audit_logs(action);
CREATE INDEX idx_audit_created_at ON audit_logs(created_at);

CREATE TABLE medical_configurations (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    category VARCHAR(100) NOT NULL,
    config_key VARCHAR(150) NOT NULL,
    config_val TEXT NOT NULL,
    status VARCHAR(50) DEFAULT 'active' CHECK (status IN ('active', 'inactive')),
    clinic_id UUID,
    created_by UUID,
    updated_by UUID,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (clinic_id) REFERENCES clinics(id) ON DELETE SET NULL,
    FOREIGN KEY (created_by) REFERENCES users(id) ON DELETE SET NULL,
    FOREIGN KEY (updated_by) REFERENCES users(id) ON DELETE SET NULL
);

CREATE INDEX idx_medical_config_clinic ON medical_configurations(clinic_id);
CREATE INDEX idx_medical_config_category ON medical_configurations(category, config_key);


CREATE TABLE IF NOT EXISTS clinic_staff (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID,
    clinic_id UUID,
    role_type VARCHAR(20) CHECK (role_type IN ('staff','admin')),
    status VARCHAR(20) DEFAULT 'active' CHECK (status IN ('active', 'inactive', 'suspended')),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (clinic_id) REFERENCES clinics(id) ON DELETE CASCADE
);

CREATE INDEX idx_clinic_staff_clinic ON clinic_staff(clinic_id);

CREATE TABLE IF NOT EXISTS roles (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(50) UNIQUE NOT NULL,
    description TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);


CREATE TABLE IF NOT EXISTS user_roles (
    user_id UUID,
    role_id UUID,
    PRIMARY KEY(user_id, role_id),
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (role_id) REFERENCES roles(id) ON DELETE CASCADE
);
