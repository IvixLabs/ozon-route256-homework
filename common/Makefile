BIN_PATH=$(CURDIR)/bin


bin-deps:
	GOBIN=$(BIN_PATH) go install github.com/gojuno/minimock/v3/cmd/minimock@latest ; \


generate-mocks: bin-deps
	$(BIN_PATH)/minimock -i ./pkg/manager/transaction.Controller -o ./pkg/manager/transaction/mock/controller_mock.go -n ControllerMock ; \
	$(BIN_PATH)/minimock -i ./pkg/manager/transaction.Transaction -o ./pkg/manager/transaction/mock/transaction_mock.go -n TransactionMock
