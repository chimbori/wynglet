-- name: InsertRating :exec
INSERT INTO ratings (url, ui, rating, ip_address, rated_at)
  VALUES ($1, $2, $3, $4, NOW());

-- name: HasRecentRatingByIPForURL :one
SELECT EXISTS (
  SELECT 1 FROM ratings
  WHERE url = $1
    AND ip_address = $2
    AND rated_at >= NOW() - INTERVAL '24 hours'
);

-- name: CountRatingGroups :one
SELECT COUNT(*) FROM (
  SELECT 1 FROM ratings GROUP BY url, ui
) AS grouped;

-- name: ListRatingGroupsPaginated :many
SELECT
    url,
    ui,
    COUNT(*)::bigint AS total_ratings,
    COALESCE(SUM(CASE WHEN rating = 1 THEN 1 ELSE 0 END), 0)::bigint AS thumbs_up,
    COALESCE(SUM(CASE WHEN rating = -1 THEN 1 ELSE 0 END), 0)::bigint AS thumbs_down,
    COALESCE(AVG(CASE WHEN ui = 'stars' THEN rating::float8 ELSE NULL END), 0)::float8 AS average_stars,
    MAX(rated_at) AS last_rated_at
  FROM ratings
  GROUP BY url, ui
  ORDER BY last_rated_at DESC
  LIMIT $1 OFFSET $2;

-- name: GetRatingsByURL :many
SELECT
    url,
    ui,
    COUNT(*)::bigint AS total_ratings,
    COALESCE(SUM(CASE WHEN rating = 1 THEN 1 ELSE 0 END), 0)::bigint AS thumbs_up,
    COALESCE(SUM(CASE WHEN rating = -1 THEN 1 ELSE 0 END), 0)::bigint AS thumbs_down,
    COALESCE(AVG(CASE WHEN ui = 'stars' THEN rating::float8 ELSE NULL END), 0)::float8 AS average_stars
  FROM ratings
  GROUP BY url, ui
  ORDER BY total_ratings DESC, url ASC;

-- name: GetRatingsByDay :many
SELECT
    to_char(date_trunc('day', rated_at), 'YYYY-MM-DD') AS day,
    ui,
    COUNT(*)::bigint AS total_ratings,
    COALESCE(SUM(CASE WHEN rating = 1 THEN 1 ELSE 0 END), 0)::bigint AS thumbs_up,
    COALESCE(SUM(CASE WHEN rating = -1 THEN 1 ELSE 0 END), 0)::bigint AS thumbs_down,
    COALESCE(AVG(CASE WHEN ui = 'stars' THEN rating::float8 ELSE NULL END), 0)::float8 AS average_stars
  FROM ratings
  WHERE rated_at >= NOW() - ($1 * INTERVAL '1 day')
  GROUP BY day, ui
  ORDER BY day DESC, ui ASC;

-- name: DeleteOldRatings :execrows
DELETE FROM ratings
  WHERE rated_at < NOW() - $1::interval;
