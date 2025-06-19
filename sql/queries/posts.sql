-- name: CreatePost :one
INSERT INTO posts (title, url, feed_id)
VALUES (
    $1,
    $2,
    $3
)

RETURNING *;

-- name: UpdateDescription :one
UPDATE posts
    SET description = $1, updated_at = NOW()
    WHERE id = $2
RETURNING *;


-- name: UpdatePublishedAt :one
UPDATE posts
    SET published_at = $1, updated_at = NOW()
    WHERE id = $2
RETURNING *;

-- name: GetPosts :many
SELECT posts.*
    FROM posts
    JOIN feed_follows ON posts.feed_id = feed_follows.feed_id
        WHERE feed_follows.user_id = $1
    ORDER BY posts.published_at DESC NULLS LAST
    LIMIT $2;