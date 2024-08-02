BIN_PATH=$(CURDIR)/bin

bin-deps:
	GOBIN=$(BIN_PATH) go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest

generate-sql: bin-deps
	${BIN_PATH}/sqlc -f sqlc/sqlc.json generate ;\
	go mod tidy
