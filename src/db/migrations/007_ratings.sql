-- +goose Up

CREATE TABLE ratings (
  _id         BIGSERIAL PRIMARY KEY,
  url         TEXT NOT NULL,
  ui          TEXT NOT NULL CHECK (ui IN ('thumbs', 'stars')),
  rating      SMALLINT NOT NULL,
  ip_address  TEXT NOT NULL,
  rated_at    TIMESTAMPTZ NOT NULL DEFAULT NOW(),

  CONSTRAINT chk_ratings_scale_values CHECK (
    (ui = 'thumbs' AND rating IN (-1, 1)) OR
    (ui = 'stars' AND rating BETWEEN 1 AND 5)
  )
);

CREATE INDEX idx_ratings_url_ui ON ratings(url, ui);
CREATE INDEX idx_ratings_rated_at ON ratings(rated_at DESC);
CREATE INDEX idx_ratings_url_ip_rated_at ON ratings(url, ip_address, rated_at DESC);
