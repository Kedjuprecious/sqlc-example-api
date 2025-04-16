-- name: CreateMessage :one
INSERT INTO message (thread, sender, content)
VALUES ($1, $2, $3)
RETURNING *;

-- name: GetMessageByID :one
SELECT * FROM message
WHERE id = $1;

-- name: GetMessagesByThread :many
SELECT * FROM message
WHERE thread = $1
ORDER BY created_at DESC;

-- name: DeleteMessage :exec
DELETE FROM message
WHERE id = $1;

-- name: EditMessage :one
UPDATE message
SET content = $2
WHERE id = $1
RETURNING *;

-- name: SearchMessagesByKeyword :many
SELECT * FROM message 
WHERE thread = $1 AND content ILIKE '%' || $2 || '%'
ORDER BY created_at DESC;

-- name: CreateThread :one
INSERT INTO thread (description) 
VALUES ($1)
RETURNING *;

-- name: GetThreadID :one
SELECT * FROM thread
WHERE id = $1;

-- name: CountMessagesInThread :one
SELECT COUNT(*) FROM message
WHERE thread = $1;

-- name: GetLatestMessageInThread :one
SELECT * FROM message
WHERE thread = $1
ORDER BY created_at DESC
LIMIT 1;

-- name: DeleteThread :exec
DELETE FROM thread 
WHERE id = $1;

-- name: GetMessagesByThreadPaginated :many
SELECT * FROM message
WHERE thread = $1
ORDER BY created_at ASC
LIMIT $2 OFFSET $3;
