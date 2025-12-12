-- name: CreatePost :one
INSERT INTO posts(created_at, updated_at, title, url, description, published_at, feed_id)
VALUES (
	NOW(),
	NOW(),
	$1,
	$2,
	$3,
	$4,
	$5
)
RETURNING *;

-- name: GetPostsForUser :many
SELECT *
FROM posts
WHERE feed_id in (
	SELECT feed_id
	FROM feed_follows
	WHERE user_id = $1
)
ORDER BY created_at DESC
LIMIT cast($2 AS INT);
