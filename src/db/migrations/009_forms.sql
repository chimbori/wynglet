-- +goose Up

CREATE TABLE form_submissions (
  _id             BIGSERIAL PRIMARY KEY,
  form_id         TEXT      NOT NULL,
  submitted_at    TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  domain          TEXT      NOT NULL,
  ip_address      TEXT      NOT NULL,
  form_data       TEXT      NOT NULL,  -- YAML-formatted.
  is_spam         BOOLEAN   NOT NULL DEFAULT FALSE,
  email_sent_at   TIMESTAMPTZ  -- NULL implies email has not yet been sent.
);

CREATE INDEX idx_form_submissions_form_id ON form_submissions(form_id);
CREATE INDEX idx_form_submissions_submitted_at ON form_submissions(submitted_at DESC);
CREATE INDEX idx_form_submissions_form_id_submitted_at ON form_submissions(form_id, submitted_at DESC);
