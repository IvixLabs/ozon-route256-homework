package shard

type NoopManager struct {
}

func NewNoopManager() *NoopManager {
	return &NoopManager{}
}

func (m *NoopManager) GetRandShardIndex() Index {
	return Index(0)
}
