-- +goose Up

ALTER TABLE logs ADD COLUMN ip TEXT;

CREATE INDEX idx_logs_ip ON logs(ip);
