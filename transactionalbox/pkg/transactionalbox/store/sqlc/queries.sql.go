// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0
// source: queries.sql

package sqlc

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)

const findLockedRecords = `-- name: FindLockedRecords :many
SELECT id, state, created_at, message_topic, message_key, message_body FROM outbox WHERE state=$1 ORDER BY created_at LIMIT $2 FOR UPDATE
`

type FindLockedRecordsParams struct {
	State int16
	Limit int32
}

func (q *Queries) FindLockedRecords(ctx context.Context, arg FindLockedRecordsParams) ([]Outbox, error) {
	rows, err := q.db.Query(ctx, findLockedRecords, arg.State, arg.Limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Outbox
	for rows.Next() {
		var i Outbox
		if err := rows.Scan(
			&i.ID,
			&i.State,
			&i.CreatedAt,
			&i.MessageTopic,
			&i.MessageKey,
			&i.MessageBody,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const insertRecord = `-- name: InsertRecord :exec
INSERT INTO outbox (id, state, created_at, message_topic, message_key, message_body) VALUES ($1, $2, $3, $4, $5, $6)
`

type InsertRecordParams struct {
	ID           pgtype.UUID
	State        int16
	CreatedAt    pgtype.Timestamp
	MessageTopic pgtype.Text
	MessageKey   []byte
	MessageBody  []byte
}

func (q *Queries) InsertRecord(ctx context.Context, arg InsertRecordParams) error {
	_, err := q.db.Exec(ctx, insertRecord,
		arg.ID,
		arg.State,
		arg.CreatedAt,
		arg.MessageTopic,
		arg.MessageKey,
		arg.MessageBody,
	)
	return err
}

const setState = `-- name: SetState :exec
UPDATE outbox SET state=$1 WHERE id =ANY($2::uuid[])
`

type SetStateParams struct {
	State int16
	Ids   []pgtype.UUID
}

func (q *Queries) SetState(ctx context.Context, arg SetStateParams) error {
	_, err := q.db.Exec(ctx, setState, arg.State, arg.Ids)
	return err
}
