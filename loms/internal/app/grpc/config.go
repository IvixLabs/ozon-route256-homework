package grpc

type StorageMode string

var InMemoryStorageMode = StorageMode("in_memory")
var SqlcStorageMode = StorageMode("sqlc")

type KafkaMode string

var KafkaKafkaMode = KafkaMode("kafka")
var NoopKafkaMode = StorageMode("nool")

type Config struct {
	GRPCAddr         string
	DBString         string
	DBReplicaString  string
	StorageMode      StorageMode
	KafkaMode        KafkaMode
	DBString1        string
	DBReplicaString1 string
}
