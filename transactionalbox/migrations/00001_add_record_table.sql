-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS outbox
(
  id uuid not null,
  state smallint not null,
  created_at timestamp,
  message_topic text,
  message_key bytea not null,
  message_body bytea not null,
  PRIMARY KEY (id)
);

CREATE INDEX outbox_created_at_idx ON outbox(created_at);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE outbox;
-- +goose StatementEnd
