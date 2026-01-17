-- migrations/000001_initial_schema.down.sql

DROP TABLE IF EXISTS audit_logs;
DROP TABLE IF EXISTS integration_configs;
DROP TABLE IF EXISTS time_entries;
DROP TABLE IF EXISTS invoice_items;
DROP TABLE IF EXISTS invoices;
DROP TABLE IF EXISTS clients;
DROP TABLE IF EXISTS users;
DROP EXTENSION IF EXISTS "uuid-ossp";