-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS "stock"
(
  sku int not null,
  total_count int not null,
  PRIMARY KEY (sku)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE stock;
-- +goose StatementEnd
