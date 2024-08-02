package sqlc

import (
	"context"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"route256/common/pkg/storage/sqlc"
	"strings"
	"time"
)

//type DBTX interface {
//	Exec(context.Context, string, ...interface{}) (pgconn.CommandTag, error)
//	Query(context.Context, string, ...interface{}) (pgx.Rows, error)
//	QueryRow(context.Context, string, ...interface{}) pgx.Row
//}

type Pool interface {
	sqlc.DBTX
	Begin(ctx context.Context) (pgx.Tx, error)
}

type pool struct {
	master  *pgxpool.Pool
	replica *pgxpool.Pool
}

func NewPool(master *pgxpool.Pool, replica *pgxpool.Pool) Pool {
	return &pool{
		master:  master,
		replica: replica,
	}
}

func (p *pool) Exec(ctx context.Context, sql string, arg ...interface{}) (tag pgconn.CommandTag, err error) {
	defer func(startTime time.Time) {
		recordMetrics(sql, startTime, err)
	}(time.Now())

	tag, err = p.master.Exec(ctx, sql, arg...)

	return tag, err
}

func (p *pool) Query(ctx context.Context, sql string, arg ...interface{}) (rows pgx.Rows, err error) {
	defer func(startTime time.Time) {
		recordMetrics(sql, startTime, err)
	}(time.Now())

	if isMaster(sql) {
		rows, err = p.master.Query(ctx, sql, arg...)
	} else {
		rows, err = p.replica.Query(ctx, sql, arg...)
	}

	return rows, err
}

func (p *pool) QueryRow(ctx context.Context, sql string, arg ...interface{}) (row pgx.Row) {
	defer func(startTime time.Time) {
		recordMetrics(sql, startTime, nil)
	}(time.Now())

	if isMaster(sql) {
		row = p.master.QueryRow(ctx, sql, arg...)
	} else {
		row = p.replica.QueryRow(ctx, sql, arg...)
	}

	return row
}

func (p *pool) Begin(ctx context.Context) (pgx.Tx, error) {
	return p.master.Begin(ctx)
}

func isMaster(sql string) bool {
	return strings.Index(sql, "FOR UPDATE") != -1 || strings.Index(sql, "INSERT INTO") != -1
}
