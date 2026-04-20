-- name: InsertRating :exec
INSERT INTO ratings (url, ui, rating, ip_address, rated_at)
  VALUES (@url, @ui, @rating, @ip_address, NOW());

-- name: HasRecentRatingByIPForURL :one
SELECT EXISTS (
  SELECT 1 FROM ratings
  WHERE url = @url
    AND ip_address = @ip_address
    AND rated_at >= NOW() - INTERVAL '24 hours'
);

-- name: CountRecentRatingsByIP :one
SELECT COUNT(*) FROM ratings
WHERE ip_address = @ip_address
  AND rated_at >= NOW() - INTERVAL '1 hour';

-- name: CountRatingGroups :one
SELECT COUNT(*) FROM (
  SELECT 1 FROM ratings
  WHERE (sqlc.arg(days) = 0 OR rated_at >= NOW() - (sqlc.arg(days) * INTERVAL '1 day'))
  GROUP BY url, ui
) AS grouped;

-- name: ListRatingsWithDistribution :many
SELECT
    url,
    ui,
    COUNT(*)::bigint AS total_ratings,
    COALESCE(SUM(CASE WHEN ui = 'thumbs' AND rating = 1 THEN 1 ELSE 0 END), 0)::bigint AS thumbs_up,
    COALESCE(SUM(CASE WHEN ui = 'thumbs' AND rating = -1 THEN 1 ELSE 0 END), 0)::bigint AS thumbs_down,
    COALESCE(SUM(CASE WHEN ui = 'stars' AND rating = 1 THEN 1 ELSE 0 END), 0)::bigint AS stars_1,
    COALESCE(SUM(CASE WHEN ui = 'stars' AND rating = 2 THEN 1 ELSE 0 END), 0)::bigint AS stars_2,
    COALESCE(SUM(CASE WHEN ui = 'stars' AND rating = 3 THEN 1 ELSE 0 END), 0)::bigint AS stars_3,
    COALESCE(SUM(CASE WHEN ui = 'stars' AND rating = 4 THEN 1 ELSE 0 END), 0)::bigint AS stars_4,
    COALESCE(SUM(CASE WHEN ui = 'stars' AND rating = 5 THEN 1 ELSE 0 END), 0)::bigint AS stars_5,
    COALESCE(AVG(CASE WHEN ui = 'stars' THEN rating::float8 ELSE NULL END), 0)::float8 AS average_stars,
    CASE
      WHEN ui = 'thumbs' THEN
        SUM(CASE WHEN rating = 1 THEN 1 ELSE 0 END)::float8 / NULLIF(COUNT(*), 0)::float8
      ELSE
        COALESCE(AVG(rating::float8), 0) / 5.0
    END::float8 AS normalized_score
  FROM ratings
  WHERE (sqlc.arg(days) = 0 OR rated_at >= NOW() - (sqlc.arg(days) * INTERVAL '1 day'))
  GROUP BY url, ui
  ORDER BY normalized_score DESC, url ASC
  LIMIT sqlc.arg(pagination_limit) OFFSET sqlc.arg(pagination_offset);

-- name: DeleteOldRatings :execrows
DELETE FROM ratings
  WHERE rated_at < NOW() - $1::interval;
