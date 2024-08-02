-- +goose Up
-- +goose StatementBegin
create table if not exists order_item
(
  order_id int not null,
  sku      int not null,
  count    int not null,
  PRIMARY KEY (order_id, sku)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table order_item;
-- +goose StatementEnd
