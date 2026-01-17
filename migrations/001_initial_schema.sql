-- migrations/001_initial_schema.sql

-- Enable UUID extension
CREATE
EXTENSION IF NOT EXISTS "uuid-ossp";

-- Users table
CREATE TABLE users
(
    id            UUID PRIMARY KEY         DEFAULT uuid_generate_v4(),
    email         VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255)        NOT NULL,
    full_name     VARCHAR(255)        NOT NULL,
    company_name  VARCHAR(255),
    created_at    TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at    TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Clients table (invoice recipients)
CREATE TABLE clients
(
    id           UUID PRIMARY KEY         DEFAULT uuid_generate_v4(),
    user_id      UUID         NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    name         VARCHAR(255) NOT NULL,
    email        VARCHAR(255),
    company_name VARCHAR(255),
    address      TEXT,
    phone        VARCHAR(50),
    created_at   TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at   TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Invoices table
CREATE TABLE invoices
(
    id                UUID PRIMARY KEY         DEFAULT uuid_generate_v4(),
    user_id           UUID               NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    client_id         UUID               NOT NULL REFERENCES clients (id) ON DELETE RESTRICT,
    invoice_number    VARCHAR(50) UNIQUE NOT NULL,
    status            VARCHAR(20)        NOT NULL CHECK (status IN ('draft', 'sent', 'paid', 'overdue', 'cancelled')),
    issue_date        DATE               NOT NULL,
    due_date          DATE               NOT NULL,
    subtotal          DECIMAL(12, 2)     NOT NULL,
    tax_rate          DECIMAL(5, 2)            DEFAULT 0,
    tax_amount        DECIMAL(12, 2)           DEFAULT 0,
    total             DECIMAL(12, 2)     NOT NULL,
    currency          VARCHAR(3)               DEFAULT 'USD',
    notes             TEXT,

    -- Integration fields
    square_invoice_id VARCHAR(255),
    square_payment_id VARCHAR(255),

    created_at        TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at        TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,

    INDEX             idx_user_invoices (user_id, created_at DESC),
    INDEX             idx_client_invoices (client_id),
    INDEX             idx_invoice_status (status),
    INDEX             idx_square_invoice (square_invoice_id)
);

-- Invoice items table
CREATE TABLE invoice_items
(
    id          UUID PRIMARY KEY         DEFAULT uuid_generate_v4(),
    invoice_id  UUID           NOT NULL REFERENCES invoices (id) ON DELETE CASCADE,
    description TEXT           NOT NULL,
    quantity    DECIMAL(10, 2) NOT NULL,
    unit_price  DECIMAL(12, 2) NOT NULL,
    amount      DECIMAL(12, 2) NOT NULL,
    sort_order  INT                      DEFAULT 0,
    created_at  TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Time entries table
CREATE TABLE time_entries
(
    id              UUID PRIMARY KEY         DEFAULT uuid_generate_v4(),
    user_id         UUID          NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    invoice_id      UUID          REFERENCES invoices (id) ON DELETE SET NULL,
    description     TEXT          NOT NULL,
    hours           DECIMAL(5, 2) NOT NULL,
    hourly_rate     DECIMAL(12, 2),
    date            DATE          NOT NULL,

    -- Jira integration
    jira_issue_key  VARCHAR(50),
    jira_worklog_id VARCHAR(50),
    jira_synced_at  TIMESTAMP WITH TIME ZONE,

    is_billable     BOOLEAN                  DEFAULT true,
    is_invoiced     BOOLEAN                  DEFAULT false,

    created_at      TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at      TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,

    INDEX           idx_user_entries (user_id, date DESC),
    INDEX           idx_invoice_entries (invoice_id),
    INDEX           idx_jira_worklog (jira_worklog_id),
    INDEX           idx_uninvoiced_billable (user_id, is_billable, is_invoiced)
);

-- Integration configurations (encrypted sensitive data)
CREATE TABLE integration_configs
(
    id          UUID PRIMARY KEY         DEFAULT uuid_generate_v4(),
    user_id     UUID        NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    provider    VARCHAR(50) NOT NULL, -- 'jira', 'square'
    config_data JSONB       NOT NULL, -- encrypted credentials
    is_active   BOOLEAN                  DEFAULT true,
    created_at  TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at  TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,

    UNIQUE (user_id, provider)
);

-- Audit log for compliance
CREATE TABLE audit_logs
(
    id          UUID PRIMARY KEY         DEFAULT uuid_generate_v4(),
    user_id     UUID        REFERENCES users (id) ON DELETE SET NULL,
    entity_type VARCHAR(50) NOT NULL,
    entity_id   UUID        NOT NULL,
    action      VARCHAR(50) NOT NULL,
    changes     JSONB,
    ip_address  INET,
    created_at  TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,

    INDEX       idx_audit_entity (entity_type, entity_id),
    INDEX       idx_audit_user (user_id, created_at DESC)
);