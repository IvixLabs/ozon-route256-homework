-- +goose Up
-- +goose StatementBegin
TRUNCATE stock;

INSERT INTO stock(sku, total_count) VALUES
                                      (1076963, 4),
                                      (1148162, 5),
                                      (1625903, 3),
                                      (2618151, 2),
                                      (2956315, 1),
                                      (2958025, 3),
                                      (3596599, 2);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE stock;
-- +goose StatementEnd


