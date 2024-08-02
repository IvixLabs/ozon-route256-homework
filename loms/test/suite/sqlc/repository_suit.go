package sqlc

import (
	"context"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
	"os"
	"route256/loms/internal/logger"
	"route256/loms/internal/manager/shard"
	"route256/loms/internal/model"
	"route256/loms/internal/repository/order"
	"route256/loms/internal/repository/stock"
	"route256/loms/internal/storage/sqlc"
	order2 "route256/loms/internal/usecase/order"
	stock2 "route256/loms/internal/usecase/stock"
	"testing"
	"time"
)

type RepositorySuite struct {
	suite.Suite
	ctx          context.Context
	postgresCont *postgres.PostgresContainer
	dbString     string
	pool         sqlc.Pool
}

func (s *RepositorySuite) SetupSuite() {
	s.ctx = context.Background()

	dbName := "loms"
	dbUser := "user"
	dbPassword := "password"

	var err error
	s.postgresCont, err = postgres.RunContainer(s.ctx,
		testcontainers.WithImage("docker.io/postgres:16-alpine"),
		postgres.WithDatabase(dbName),
		postgres.WithUsername(dbUser),
		postgres.WithPassword(dbPassword),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(120*time.Second)),
	)

	if err != nil {
		panic(err)
	}

	s.dbString, err = s.postgresCont.ConnectionString(s.ctx, "sslmode=disable")

}

func (s *RepositorySuite) TearDownSuite() {
	_ = s.postgresCont.Terminate(s.ctx)

}

func (s *RepositorySuite) SetupTest() {
	var err error

	db, err := goose.OpenDBWithDriver("postgres", s.dbString)
	if err != nil {
		panic(err)
	}

	err = goose.UpContext(s.ctx, db, os.Getenv("MIGRATIONS_PATH"))
	if err != nil {
		panic(err)
	}

	err = goose.UpContext(s.ctx, db, os.Getenv("TEST_MIGRATIONS_PATH"))
	if err != nil {
		panic(err)
	}

	s.pool, err = sqlc.GetPool(s.ctx, s.dbString, "")
	if err != nil {
		panic(err)
	}
}

func (s *RepositorySuite) TearDownTest() {
	db, err := goose.OpenDBWithDriver("postgres", s.dbString)
	if err != nil {
		panic(err)
	}

	err = goose.DownContext(s.ctx, db, os.Getenv("TEST_MIGRATIONS_PATH"))
	if err != nil {
		panic(err)
	}

	err = goose.DownContext(s.ctx, db, os.Getenv("MIGRATIONS_PATH"))
	if err != nil {
		panic(err)
	}
}

func (s *RepositorySuite) TestStocks() {

	txController := sqlc.NewTransactionController(s.pool, 0)

	stockRepo := stock.NewSqlcRepository(txController, logger.NewStubLogger())

	tests := []struct {
		name      string
		operation func(*testing.T)
	}{
		{
			name: "Get_stock",
			operation: func(t *testing.T) {
				stockObj, err := stockRepo.GetBySku(s.ctx, 1076963)
				assert.NotNil(t, stockObj)
				assert.NoError(t, err)
			},
		},
		{
			name: "Get_stock_not_found",
			operation: func(t *testing.T) {
				stockObj, err := stockRepo.GetBySku(s.ctx, 111)
				assert.Nil(t, stockObj)
				assert.ErrorIs(t, err, stock2.ErrStockNotFound)
			},
		},
		{
			name: "Save_stock",
			operation: func(t *testing.T) {
				stockObj := model.NewStock(111, 9)
				_, err := stockObj.Reserve(1, 2)
				assert.NoError(t, err)

				err = stockRepo.Save(s.ctx, stockObj)
				assert.NoError(t, err)
			},
		},
		{
			name: "Get_new_stock",
			operation: func(t *testing.T) {
				stockObj, err := stockRepo.GetBySku(s.ctx, 111)
				expectedStock := &model.Stock{
					Sku:        111,
					TotalCount: 7,
				}
				assert.Equal(t, expectedStock, stockObj)
				assert.NoError(t, err)
			},
		},
	}

	for _, tt := range tests {
		s.T().Run(tt.name, func(t *testing.T) {
			tt.operation(t)
		})
	}
}

func (s *RepositorySuite) TestOrders() {

	txController := sqlc.NewTransactionController(s.pool, 0)
	sqlcTxs := []*sqlc.TransactionController{txController}
	shardManager := shard.NewSqlcManager(sqlcTxs)

	orderRepo := order.NewSqlcRepository(shardManager, logger.NewStubLogger())

	var newOrderId model.OrderID
	tests := []struct {
		name      string
		operation func(*testing.T)
	}{
		{
			name: "Get_order_not_found",
			operation: func(t *testing.T) {
				orderObj, err := orderRepo.GetByID(s.ctx, 1)
				assert.Nil(t, orderObj)
				assert.ErrorIs(t, err, order2.ErrOrderNotFound)
			},
		},
		{
			name: "Save_order",
			operation: func(t *testing.T) {
				stockObj := model.NewOrder(
					111,
					[]model.OrderItem{
						{Sku: 1, Count: 2},
						{Sku: 2, Count: 3},
					},
				)
				newOrder, err := orderRepo.Save(s.ctx, stockObj)
				assert.NoError(t, err)
				assert.Greater(t, newOrder.ID, model.OrderID(0))
				newOrderId = newOrder.ID
			},
		},
		{
			name: "Get_new_order",
			operation: func(t *testing.T) {
				orderObj, err := orderRepo.GetByID(s.ctx, newOrderId)
				expectedOrder := &model.Order{
					ID:     newOrderId,
					Status: model.OrderStatusNew,
					UserID: 111,
					Items: []model.OrderItem{
						{Sku: 1, Count: 2},
						{Sku: 2, Count: 3},
					},
				}
				assert.Equal(t, expectedOrder, orderObj)
				assert.NoError(t, err)
			},
		},
	}

	for _, tt := range tests {
		s.T().Run(tt.name, func(t *testing.T) {
			tt.operation(t)
		})
	}
}
