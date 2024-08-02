-- +goose Up
-- +goose StatementBegin
CREATE TYPE order_status AS ENUM ('new', 'awaiting_payment', 'failed', 'payed', 'cancelled');

CREATE TABLE IF NOT EXISTS "order"
(
  id      serial primary key,
  user_id int          not null,
  status  order_status not null
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table "order";
-- +goose StatementEnd
