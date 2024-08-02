-- name: Save :exec
INSERT INTO reserved_stock (order_id, sku, count, status)
VALUES ($1, $2, $3, $4)
ON CONFLICT (order_id, sku) DO UPDATE
  SET count = excluded.count, status = excluded.status;

-- name: GetLocked :one
SELECT * FROM "reserved_stock" WHERE order_id = $1 AND sku = $2 FOR UPDATE;
