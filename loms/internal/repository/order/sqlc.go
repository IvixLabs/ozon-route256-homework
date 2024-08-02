package order

import (
	"context"
	"route256/loms/internal/logger"
	"route256/loms/internal/manager/shard"
	"route256/loms/internal/model"
	"route256/loms/internal/repository/order/sqlc"
	"route256/loms/internal/usecase/order"
)

type SqlcRepository struct {
	shardManager *shard.SqlcManager
	logger       logger.Logger
}

func NewSqlcRepository(
	shardManager *shard.SqlcManager,
	logger logger.Logger) *SqlcRepository {
	return &SqlcRepository{
		shardManager: shardManager,
		logger:       logger,
	}
}

func (r *SqlcRepository) Save(ctx context.Context, order *model.Order) (*model.Order, error) {
	ctx, span := beginSpan(ctx, "sqlc", "Save")
	defer span.End()

	var shardIdx shard.Index
	if order.ID > 0 {
		shardIdx = r.shardManager.GetShardIndexFromID(int64(order.ID))
	} else {
		shardIdx = r.shardManager.GetRandShardIndex()
	}

	tc, pickErr := r.shardManager.Pick(shardIdx)

	if pickErr != nil {
		return nil, pickErr
	}

	queries := sqlc.New(tc.GetDBTX(ctx))

	newOrder := *order
	if order.ID > 0 {
		err := queries.UpdateOrder(ctx, sqlc.UpdateOrderParams{
			ID:     int32(order.ID),
			UserID: int32(order.UserID),
			Status: sqlc.OrderStatus(order.Status),
		})
		if err != nil {
			r.logger.Error(err)
			return nil, err
		}
	} else {
		orderId, err := queries.InsertOrder(ctx, sqlc.InsertOrderParams{
			ShardID: int32(shardIdx),
			UserID:  int32(order.UserID),
			Status:  sqlc.OrderStatus(order.Status),
		})
		if err != nil {
			r.logger.Error(err)
			return nil, err
		}
		newOrder.ID = model.OrderID(orderId)
	}

	for _, orderItem := range newOrder.Items {
		err := queries.SaveOrderItem(ctx, sqlc.SaveOrderItemParams{
			OrderID: int32(newOrder.ID),
			Sku:     int32(orderItem.Sku),
			Count:   int32(orderItem.Count),
		})
		if err != nil {
			r.logger.Error(err)
			return nil, err
		}
	}

	return &newOrder, nil
}

func (r *SqlcRepository) GetByID(ctx context.Context, orderId model.OrderID) (*model.Order, error) {
	ctx, span := beginSpan(ctx, "sqlc", "GetByID")
	defer span.End()

	shardIdx := r.shardManager.GetShardIndexFromID(int64(orderId))
	tc, err := r.shardManager.Pick(shardIdx)
	if err != nil {
		return nil, err
	}

	queries := sqlc.New(tc.GetDBTX(ctx))

	orderParts, err := queries.GetByID(ctx, int32(orderId))
	if err != nil {
		r.logger.Error(err)
		return nil, err
	}
	if len(orderParts) == 0 {
		return nil, order.ErrOrderNotFound
	}

	firstRawOrderItem := orderParts[0]
	orderEntity := &model.Order{
		ID:     model.OrderID(firstRawOrderItem.ID),
		Status: model.OrderStatus(firstRawOrderItem.Status),
		UserID: model.UserID(firstRawOrderItem.UserID),
		Items:  make([]model.OrderItem, len(orderParts)),
	}

	for i, rawOrderItem := range orderParts {
		orderEntity.Items[i] = model.OrderItem{
			Sku:   model.Sku(rawOrderItem.Sku.Int32),
			Count: model.Count(rawOrderItem.Count.Int32),
		}
	}

	return orderEntity, nil
}

func (r *SqlcRepository) GetLockByID(ctx context.Context, orderId model.OrderID) (*model.Order, error) {
	ctx, span := beginSpan(ctx, "sqlc", "GetLockByID")
	defer span.End()

	shardIdx := r.shardManager.GetShardIndexFromID(int64(orderId))
	tc, err := r.shardManager.Pick(shardIdx)
	if err != nil {
		return nil, err
	}

	queries := sqlc.New(tc.GetDBTX(ctx))

	orderParts, err := queries.GetLockByID(ctx, int32(orderId))
	if err != nil {
		r.logger.Error(err)
		return nil, err
	}
	if len(orderParts) == 0 {
		return nil, order.ErrOrderNotFound
	}

	firstRawOrderItem := orderParts[0]
	orderEntity := &model.Order{
		ID:     model.OrderID(firstRawOrderItem.ID),
		Status: model.OrderStatus(firstRawOrderItem.Status),
		UserID: model.UserID(firstRawOrderItem.UserID),
		Items:  make([]model.OrderItem, len(orderParts)),
	}

	for i, rawOrderItem := range orderParts {
		orderEntity.Items[i] = model.OrderItem{
			Sku:   model.Sku(rawOrderItem.Sku.Int32),
			Count: model.Count(rawOrderItem.Count.Int32),
		}
	}

	return orderEntity, nil
}
