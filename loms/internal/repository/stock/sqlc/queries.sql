-- name: GetBySku :one
SELECT * FROM stock WHERE stock.sku = @stock_sku;

-- name: GetLockedBySku :one
SELECT * FROM stock WHERE stock.sku = @stock_sku FOR UPDATE OF stock;

-- name: Save :exec
INSERT INTO stock (sku, total_count) VALUES ($1, $2)
ON CONFLICT (sku) DO UPDATE SET total_count = excluded.total_count;
