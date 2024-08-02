LOMS_PROTO_PATH=loms/v1
MESSAGE_PROTO_PATH=message/v1

BIN_PATH=$(CURDIR)/bin

APP_PG_PASS=master

TEST_DIRS=./internal/usecase/... ./internal/repository/... ./internal/model/... ./internal/manager...
bin-deps:
	GOBIN=$(BIN_PATH) go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28.1 ; \
	GOBIN=$(BIN_PATH) go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2.0 ; \
	GOBIN=$(BIN_PATH) go install github.com/bufbuild/buf/cmd/buf@v1.21.0 ; \
	GOBIN=$(BIN_PATH) go install github.com/envoyproxy/protoc-gen-validate@v1.0.4 ; \
	GOBIN=$(BIN_PATH) go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway@v2.19.1 ;\
	GOBIN=$(BIN_PATH) go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2@v2.19.1 ;\
	GOBIN=$(BIN_PATH) go install github.com/gojuno/minimock/v3/cmd/minimock@latest ; \
	GOBIN=$(BIN_PATH) go install github.com/fzipp/gocyclo/cmd/gocyclo@latest ; \
	GOBIN=$(BIN_PATH) go install github.com/uudashr/gocognit/cmd/gocognit@latest ; \
	GOBIN=$(BIN_PATH) go install gotest.tools/gotestsum@latest ; \
	GOBIN=$(BIN_PATH) go install github.com/pressly/goose/v3/cmd/goose@latest ; \
	GOBIN=$(BIN_PATH) go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest ; \
	GOBIN=$(BIN_PATH) go install github.com/boumenot/gocover-cobertura@latest

.PHONY: build
build:
	go build -o dist/app cmd/app/*

run : build
	set -a && source .env && set +a && dist/app_grpc

generate-api: bin-deps vendor-proto
	protoc \
	-I api/${LOMS_PROTO_PATH} \
	-I vendor-proto \
	--plugin=protoc-gen-go=${BIN_PATH}/protoc-gen-go \
	--go_out internal/pb/${LOMS_PROTO_PATH} \
	--go_opt paths=source_relative \
	--plugin=protoc-gen-go-grpc=${BIN_PATH}/protoc-gen-go-grpc \
	--go-grpc_out internal/pb/${LOMS_PROTO_PATH} \
	--go-grpc_opt paths=source_relative \
	--plugin=protoc-gen-validate=${BIN_PATH}/protoc-gen-validate \
	--validate_out="lang=go,paths=source_relative:internal/pb/${LOMS_PROTO_PATH}" \
	--plugin=protoc-gen-grpc-gateway=${BIN_PATH}/protoc-gen-grpc-gateway \
	--grpc-gateway_out internal/pb/${LOMS_PROTO_PATH} \
	--grpc-gateway_opt logtostderr=true \
	--grpc-gateway_opt paths=source_relative \
	--grpc-gateway_opt generate_unbound_methods=true \
	--plugin=protoc-gen-openapiv2=${BIN_PATH}/protoc-gen-openapiv2 \
	--openapiv2_out internal/swagger/assets \
	--openapiv2_opt logtostderr=true \
	api/${LOMS_PROTO_PATH}/loms.proto ; \
	protoc \
    -I api/${LOMS_PROTO_PATH} \
    -I vendor-proto \
	--plugin=protoc-gen-go=${BIN_PATH}/protoc-gen-go \
	--go_out internal/pb/${LOMS_PROTO_PATH} \
	--go_opt paths=source_relative \
	api/${LOMS_PROTO_PATH}/message.proto ; \
	go mod tidy

run-unit-tests:
	go test -p 10 -race -v ${TEST_DIRS}

run-unit-tests-coverage:
	mkdir reports ; \
	go test -p 10 -race -coverprofile=reports/coverage.tmp -v ${TEST_DIRS} && \
	cat reports/coverage.tmp | grep -v '_mock.go' > reports/coverage.out && \
    go tool cover -html=reports/coverage.out -o reports/coverage.html

run-e2e-tests:
	DOCKER_CONTEXT_PATH=$(CURDIR) go test -v ./test/cmd/e2e

run-sqlc-tests:
	DOCKER_CONTEXT_PATH=$(CURDIR) \
	MIGRATIONS_PATH=$(CURDIR)/migrations \
	TEST_MIGRATIONS_PATH=$(CURDIR)/test_migrations \
 	go test -v ./test/cmd/sqlc


ci-run-unit-tests:
	mkdir reports ; \
	${BIN_PATH}/gotestsum --junitfile ./reports/report.xml --format testname -- -coverprofile=reports/coverage.out ${TEST_DIRS} ; \
	$(BIN_PATH)/gocover-cobertura < reports/coverage.out > reports/coverage.xml

run-gocyclo-linter:
	$(BIN_PATH)/gocyclo -ignore "_test|_mock|\.pb\." -over 12  ./internal ./cmd

run-gocognit-linter:
	$(BIN_PATH)/gocognit -ignore "_test|_mock|\.pb\." -over 22  ./internal ./cmd

generate-mocks: bin-deps
	$(BIN_PATH)/minimock -i ./internal/usecase/order.Repository -o ./internal/usecase/order/mock/repository_mock.go -n RepositoryMock ; \
	$(BIN_PATH)/minimock -i ./internal/usecase/order.StockService -o ./internal/usecase/order/mock/stock_service_mock.go -n StockServiceMock ; \
	$(BIN_PATH)/minimock -i ./internal/usecase/stock.Repository -o ./internal/usecase/stock/mock/repository_mock.go -n RepositoryMock ;\
	$(BIN_PATH)/minimock -i ./internal/usecase/stock.ReservedStockRepository -o ./internal/usecase/stock/mock/reserved_stock_repository_mock.go -n ReservedStockRepositoryMock ;\
	$(BIN_PATH)/minimock -i ./internal/usecase/order.Publisher -o ./internal/usecase/order/mock/publisher_mock.go -n PublisherMock

vendor-proto/validate:
	git clone -b main --single-branch --depth=2 --filter=tree:0 \
		https://github.com/bufbuild/protoc-gen-validate vendor-proto/tmp && \
		cd vendor-proto/tmp && \
		git sparse-checkout set --no-cone validate &&\
		git checkout
		mkdir -p vendor-proto/validate
		mv vendor-proto/tmp/validate vendor-proto/
		rm -rf vendor-proto/tmp

vendor-proto/protoc-gen-openapiv2/options:
	git clone -b main --single-branch -n --depth=1 --filter=tree:0 \
 		https://github.com/grpc-ecosystem/grpc-gateway vendor-proto/grpc-ecosystem && \
 	cd vendor-proto/grpc-ecosystem && \
	git sparse-checkout set --no-cone protoc-gen-openapiv2/options && \
	git checkout
	mkdir -p vendor-proto/protoc-gen-openapiv2
	mv vendor-proto/grpc-ecosystem/protoc-gen-openapiv2/options vendor-proto/protoc-gen-openapiv2
	rm -rf vendor-proto/grpc-ecosystem

vendor-proto/google/api:
	git clone -b master --single-branch -n --depth=1 --filter=tree:0 \
 		https://github.com/googleapis/googleapis vendor-proto/googleapis && \
 	cd vendor-proto/googleapis && \
	git sparse-checkout set --no-cone google/api && \
	git checkout
	mkdir -p  vendor-proto/google
	mv vendor-proto/googleapis/google/api vendor-proto/google
	rm -rf vendor-proto/googleapis

vendor-proto: vendor-proto/validate vendor-proto/google/api vendor-proto/protoc-gen-openapiv2/options

migrate:
	set -a && source .env && set +a && ${BIN_PATH}/goose -dir migrations up ;\
    set -a && source .env && set +a && ${BIN_PATH}/goose -dir ../transactionalbox/migrations -table transactionalbox_goose_db_version up

generate-sql:
	${BIN_PATH}/sqlc -f sqlc/sqlc.json generate
