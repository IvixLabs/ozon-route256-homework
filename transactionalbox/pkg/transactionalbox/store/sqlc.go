package store

import (
	"context"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"route256/common/pkg/manager/transaction"
	storageSqlc "route256/common/pkg/storage/sqlc"
	"route256/transactionalbox/pkg/transactionalbox"
	"route256/transactionalbox/pkg/transactionalbox/store/sqlc"
	"time"
)

type Sqlc struct {
	txCtrl storageSqlc.TransactionController
}

func NewSqlc(txCtrl storageSqlc.TransactionController) *Sqlc {
	return &Sqlc{txCtrl: txCtrl}
}

var _ transactionalbox.ProducerStore = (*Sqlc)(nil)

func (s *Sqlc) AddRecord(ctx context.Context, record transactionalbox.Record) error {

	queries := sqlc.New(s.txCtrl.GetDBTX(ctx))

	msg := record.Message

	params := sqlc.InsertRecordParams{
		ID:           pgtype.UUID{Bytes: record.ID, Valid: true},
		State:        int16(record.State),
		CreatedAt:    pgtype.Timestamp{Time: time.Now(), Valid: true},
		MessageTopic: pgtype.Text{String: msg.Topic, Valid: true},
		MessageKey:   msg.Key,
		MessageBody:  msg.Body,
	}

	err := queries.InsertRecord(ctx, params)
	if err != nil {
		return err
	}

	return nil
}

func (s *Sqlc) FindLockedPendingRecords(ctx context.Context, size int) (records []transactionalbox.Record, retErr error) {

	tx, txCtx, err := s.txCtrl.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer func(tx transaction.Transaction, ctx context.Context) {
		commErr := tx.Commit(ctx)
		if commErr != nil {
			records = nil
			retErr = commErr
		}
	}(tx, ctx)

	qtx := sqlc.New(s.txCtrl.GetDBTX(txCtx))
	params := sqlc.FindLockedRecordsParams{
		State: int16(transactionalbox.PendingState),
		Limit: int32(size),
	}

	res, findErr := qtx.FindLockedRecords(ctx, params)
	if findErr != nil {
		return nil, findErr
	}

	lenRes := len(res)

	ids := make([]pgtype.UUID, lenRes)
	records = make([]transactionalbox.Record, lenRes)
	for i := 0; i < lenRes; i++ {
		resItem := res[i]
		ids[i] = resItem.ID

		records[i] = transactionalbox.Record{
			ID:        resItem.ID.Bytes,
			State:     transactionalbox.RecordState(resItem.State),
			CreatedAt: resItem.CreatedAt.Time,
			Message: transactionalbox.Message{
				Topic: resItem.MessageTopic.String,
				Key:   resItem.MessageKey,
				Body:  resItem.MessageBody,
			},
		}
	}

	setParams := sqlc.SetStateParams{
		State: int16(transactionalbox.ProcessingState),
		Ids:   ids,
	}
	err = qtx.SetState(ctx, setParams)
	if err != nil {
		return nil, err
	}

	return records, nil
}

func (s *Sqlc) SetDeliveredState(ctx context.Context, id uuid.UUID) error {
	queries := sqlc.New(s.txCtrl.GetDBTX(ctx))

	setParams := sqlc.SetStateParams{
		State: int16(transactionalbox.ProcessingState),
		Ids:   []pgtype.UUID{{Bytes: id, Valid: true}},
	}

	err := queries.SetState(ctx, setParams)
	if err != nil {
		return err
	}

	return nil
}

func (s *Sqlc) SetErrorState(ctx context.Context, id uuid.UUID) error {
	queries := sqlc.New(s.txCtrl.GetDBTX(ctx))

	setParams := sqlc.SetStateParams{
		State: int16(transactionalbox.ErrorState),
		Ids:   []pgtype.UUID{{Bytes: id, Valid: true}},
	}

	err := queries.SetState(ctx, setParams)
	if err != nil {
		return err
	}

	return nil
}

func (s *Sqlc) SetPendingState(ctx context.Context, id uuid.UUID) error {
	queries := sqlc.New(s.txCtrl.GetDBTX(ctx))

	setParams := sqlc.SetStateParams{
		State: int16(transactionalbox.PendingState),
		Ids:   []pgtype.UUID{{Bytes: id, Valid: true}},
	}

	err := queries.SetState(ctx, setParams)
	if err != nil {
		return err
	}

	return nil
}
