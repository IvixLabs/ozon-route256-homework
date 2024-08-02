package stock

import (
	"context"
	"errors"
	"github.com/jackc/pgx/v5"
	"route256/loms/internal/logger"
	"route256/loms/internal/model"
	"route256/loms/internal/repository/stock/sqlc"
	storageSqlc "route256/loms/internal/storage/sqlc"
	"route256/loms/internal/usecase/stock"
)

type SqlcRepository struct {
	logger       logger.Logger
	txController *storageSqlc.TransactionController
}

func NewSqlcRepository(txController *storageSqlc.TransactionController, logger logger.Logger) *SqlcRepository {
	return &SqlcRepository{
		logger:       logger,
		txController: txController,
	}
}

func (r *SqlcRepository) getQueries(ctx context.Context) *sqlc.Queries {
	return sqlc.New(r.txController.GetDBTX(ctx))
}

func (r *SqlcRepository) GetBySku(ctx context.Context, sku model.Sku) (*model.Stock, error) {
	ctx, span := beginSpan(ctx, "sqlc", "GetBySku")
	defer span.End()

	queries := r.getQueries(ctx)

	rawStock, err := queries.GetBySku(ctx, int32(sku))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, stock.ErrStockNotFound
		}

		r.logger.Error(err)
		return nil, err
	}

	stockObj := &model.Stock{
		Sku:        model.Sku(rawStock.Sku),
		TotalCount: model.Count(rawStock.TotalCount),
	}

	return stockObj, nil
}

func (r *SqlcRepository) Save(ctx context.Context, stockObj *model.Stock) error {
	ctx, span := beginSpan(ctx, "sqlc", "Save")
	defer span.End()

	queries := r.getQueries(ctx)

	err := queries.Save(ctx, sqlc.SaveParams{Sku: int32(stockObj.Sku), TotalCount: int32(stockObj.TotalCount)})
	if err != nil {
		r.logger.Error(err)
		return err
	}

	return nil
}

func (r *SqlcRepository) GetLockBySku(ctx context.Context, sku model.Sku) (*model.Stock, error) {
	ctx, span := beginSpan(ctx, "sqlc", "GetLockBySku")
	defer span.End()

	queries := r.getQueries(ctx)

	rawStock, err := queries.GetLockedBySku(ctx, int32(sku))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, stock.ErrStockNotFound
		}

		r.logger.Error(err)
		return nil, err
	}

	stockObj := &model.Stock{
		Sku:        model.Sku(rawStock.Sku),
		TotalCount: model.Count(rawStock.TotalCount),
	}

	return stockObj, nil
}
