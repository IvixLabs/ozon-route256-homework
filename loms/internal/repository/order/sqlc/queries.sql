-- name: GetByID :many
SELECT id, user_id, status, order_item.sku, order_item.count
FROM "order"
       LEFT JOIN order_item ON "order".id = order_item.order_id
WHERE "order".id = $1;

-- name: GetLockByID :many
SELECT id, user_id, status, order_item.sku, order_item.count
FROM "order"
       LEFT JOIN order_item ON "order".id = order_item.order_id
WHERE "order".id = $1
  FOR UPDATE OF "order";


-- name: InsertOrder :one
INSERT INTO "order" (id, user_id, status)
VALUES (nextval('order_id_manual_seq')+(@shard_id::int), $1, $2)
RETURNING id;

-- name: UpdateOrder :exec
UPDATE "order"
SET user_id = $1,
    status  = $2
WHERE id = $3;

-- name: SaveOrderItem :exec
INSERT INTO order_item (order_id, sku, count)
VALUES ($1, $2, $3)
ON CONFLICT (order_id, sku) DO UPDATE
  SET count = EXCLUDED.count;
