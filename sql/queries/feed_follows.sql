-- name: CreateFeedFollows :one
WITH new_follow as (
	INSERT INTO feed_follows(created_at, updated_at, user_id, feed_id)
	VALUES (
		$1,
		$2,
		$3,
		$4
	)
	RETURNING *
)
SELECT
	new_follow.*,
	users.name as user_name,
	feeds.name as feed_name
FROM new_follow
INNER JOIN users ON new_follow.user_id = users.id
INNER JOIN feeds ON new_follow.feed_id = feeds.id;

-- name: GetFeedFollowsForUser :many
SELECT
	feed_follows.created_at,
	feed_follows.updated_at,
	feed_follows.user_id,
	feed_follows.feed_id,
	users.name as user_name,
	feeds.name as feed_name
FROM feed_follows
INNER JOIN users on feed_follows.user_id = users.id
INNER JOIN feeds on feed_follows.feed_id = feeds.id
WHERE feed_follows.user_id = $1;

-- name: RemoveFeedFollow :exec
DELETE
FROM feed_follows
WHERE user_id = $1
AND feed_id = $2;
