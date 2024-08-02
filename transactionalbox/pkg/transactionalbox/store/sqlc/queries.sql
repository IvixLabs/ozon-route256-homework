-- name: InsertRecord :exec
INSERT INTO outbox (id, state, created_at, message_topic, message_key, message_body) VALUES ($1, $2, $3, $4, $5, $6);

-- name: FindLockedRecords :many
SELECT * FROM outbox WHERE state=$1 ORDER BY created_at LIMIT $2 FOR UPDATE ;

-- name: SetState :exec
UPDATE outbox SET state=$1 WHERE id =ANY(@ids::uuid[]);
