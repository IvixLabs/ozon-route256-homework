package reservedstock

import (
	"context"
	"errors"
	"github.com/jackc/pgx/v5"
	"route256/loms/internal/logger"
	"route256/loms/internal/model"
	"route256/loms/internal/repository/reservedstock/sqlc"
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

func (r *SqlcRepository) Save(ctx context.Context, rStock *model.ReservedStock) error {

	quqries := r.getQueries(ctx)

	params := sqlc.SaveParams{
		Sku:     int32(rStock.Sku),
		Count:   int32(rStock.Count),
		Status:  string(rStock.Status),
		OrderID: int32(rStock.OrderID),
	}
	err := quqries.Save(ctx, params)
	if err != nil {
		r.logger.Error(err)
		return err
	}

	return nil
}

func (r *SqlcRepository) GetLocked(ctx context.Context, orderID model.OrderID, sku model.Sku) (*model.ReservedStock, error) {

	quqries := r.getQueries(ctx)

	params := sqlc.GetLockedParams{
		OrderID: int32(orderID),
		Sku:     int32(sku),
	}
	raw, err := quqries.GetLocked(ctx, params)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, stock.ErrReservedStockNotFound
		}

		r.logger.Error(err)
		return nil, err
	}

	rStock := model.ReservedStock{
		Sku:     model.Sku(raw.Sku),
		OrderID: model.OrderID(raw.OrderID),
		Status:  model.ReservedStockStatus(raw.Status),
		Count:   model.Count(raw.Count),
	}

	return &rStock, nil
}
