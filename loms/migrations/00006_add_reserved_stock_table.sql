-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS "reserved_stock"
(
  order_id int not null,
  sku int not null,
  count int not null,
  status text not null,
  PRIMARY KEY (order_id, sku)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE reserved_stock;
-- +goose StatementEnd
