package shard

import (
	"errors"
	"fmt"
	"math/rand"
	"route256/loms/internal/storage/sqlc"
)

var (
	ErrShardIndexOutOfRange = errors.New("shard index is out of range")
)

type Key string
type Index int
type Fn func(Key) Index

type SqlcManager struct {
	shards      []*sqlc.TransactionController
	totalShards int
}

func NewSqlcManager(shards []*sqlc.TransactionController) *SqlcManager {
	return &SqlcManager{shards: shards, totalShards: len(shards)}
}

func (m *SqlcManager) GetRandShardIndex() Index {
	if m.totalShards == 1 {
		return 0
	}

	return Index(rand.Intn(m.totalShards))
}

func (m *SqlcManager) GetShardIndexFromID(id int64) Index {
	if m.totalShards == 1 {
		return 0
	}

	return Index(id % 1000)
}

func (m *SqlcManager) Pick(index Index) (*sqlc.TransactionController, error) {
	if int(index) < len(m.shards) {
		return m.shards[index], nil
	}

	return nil, fmt.Errorf("%w: given index=%d, len=%d", ErrShardIndexOutOfRange, index, len(m.shards))
}
