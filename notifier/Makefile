LOMS_PROTO_PATH=loms/v1

BIN_PATH=$(CURDIR)/bin

#TEST_DIRS=./internal/usecase/cart/... ./internal/repository/cart/... ./internal/model/...

.PHONY: build
build:
	go build  -o dist/app cmd/app/*

run : build
	set -a && source .env && set +a && dist/app

#generate-mocks: bin-deps
#	$(BIN_PATH)/minimock -i ./internal/usecase/cart.ProductProvider -o ./internal/usecase/cart/mock/product_provider_mock.go -n ProductProviderMock ;\
#    $(BIN_PATH)/minimock -i ./internal/usecase/cart.Repository -o ./internal/usecase/cart/mock/repository_mock.go -n RepositoryMock ;\
#    $(BIN_PATH)/minimock -i ./internal/usecase/cart.LOMSClient -o ./internal/usecase/cart/mock/loms_client_mock.go -n LOMSClientMock

#run-unit-tests:
#	go test -p 10 -race -v ${TEST_DIRS}

#run-benchmarks:
#	go test -bench=. -benchmem ./internal/repository/cart

#run-unit-tests-coverage:
#	mkdir reports ; \
#	go test -p 10 -race -coverprofile=reports/coverage.tmp -v ${TEST_DIRS} && \
#	cat reports/coverage.tmp | grep -v '_mock.go' > reports/coverage.out && \
#    go tool cover -html=reports/coverage.out -o reports/coverage.html

#run-e2e-tests:
#	DOCKER_CONTEXT_PATH=$(CURDIR) go test -v ./test/cmd

#run-gocyclo-linter:
#	$(BIN_PATH)/gocyclo -ignore "_test|_mock|\.pb\." -over 10  ./internal ./cmd

#run-gocognit-linter:
#	$(BIN_PATH)/gocognit -ignore "_test|_mock|\.pb\." -over 11  ./internal ./cmd

#ci-run-unit-tests:
#	mkdir reports ; \
#	${BIN_PATH}/gotestsum --format pkgname --junitfile ./reports/report.xml ${TEST_DIRS} ; \
#	go test -coverprofile=./reports/coverage.out -covermode=count ${TEST_DIRS} ; \
#	$(BIN_PATH)/gocover-cobertura < ./reports/coverage.out > ./reports/coverage.xml

bin-deps:
	GOBIN=$(BIN_PATH) go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28.1 ; \
	GOBIN=$(BIN_PATH) go install github.com/bufbuild/buf/cmd/buf@v1.21.0 ; \
	GOBIN=$(BIN_PATH) go install github.com/envoyproxy/protoc-gen-validate@v1.0.4

generate-api: bin-deps
	protoc \
    -I api/${LOMS_PROTO_PATH} \
    -I vendor-proto \
	--plugin=protoc-gen-go=${BIN_PATH}/protoc-gen-go \
	--go_out internal/pb/${LOMS_PROTO_PATH} \
	--go_opt paths=source_relative \
	api/${LOMS_PROTO_PATH}/message.proto ; \
	go mod tidy
