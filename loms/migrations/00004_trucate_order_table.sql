-- +goose Up
-- +goose StatementBegin
TRUNCATE TABLE order_item;
TRUNCATE TABLE "order";
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
-- +goose StatementEnd
