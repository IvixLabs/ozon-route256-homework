package loms

import (
	"context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"route256/loms/internal/model"
	"route256/loms/internal/pb/loms/v1"
)

func (s *Service) StockInfo(ctx context.Context, req *loms.StockInfoRequest) (*loms.StockInfoResponse, error) {

	stockObj, err := s.stockRepository.GetBySku(ctx, model.Sku(req.Sku))

	if err != nil {
		return nil, status.Errorf(codes.FailedPrecondition, "stock info: %v", err)
	}

	return repackToProto(stockObj), nil
}

func repackToProto(stockObj *model.Stock) *loms.StockInfoResponse {

	res := &loms.StockInfoResponse{Count: uint64(stockObj.TotalCount)}

	return res
}
