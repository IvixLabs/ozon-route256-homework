build-all-for-linux:
	cd cart      && GOOS=linux GOARCH=$(shell sh scripts/current-arch.sh) make build && cd .. ;\
	cd loms      && GOOS=linux GOARCH=$(shell sh scripts/current-arch.sh) make build && cd .. ;\
	cd notifier  && GOOS=linux GOARCH=$(shell sh scripts/current-arch.sh) make build && cd ..

build-all:
	cd cart      && GOOS=$(shell sh scripts/current-os.sh) GOARCH=$(shell sh scripts/current-arch.sh) make build && cd .. ;\
	cd loms      && GOOS=$(shell sh scripts/current-os.sh) GOARCH=$(shell sh scripts/current-arch.sh) make build && cd .. ;\
	cd notifier  && GOOS=$(shell sh scripts/current-os.sh) GOARCH=$(shell sh scripts/current-arch.sh) make build && cd ..

run-all: build-all
	go run launcher/cmd/launcher/main.go \
	$(CURDIR)/cart/.env     $(CURDIR)/cart/dist/app \
	$(CURDIR)/loms/.env     $(CURDIR)/loms/dist/app \
	$(CURDIR)/notifier/.env $(CURDIR)/notifier/dist/app

run-all-migrates:
	cd loms && make migrate && cd ..

run-all-dc: build-all-for-linux
	docker compose up --force-recreate --build

run-dc-db:
	docker compose up --force-recreate --build \
	loms_sync_slave_db loms_master_db loms_migrate \
	redis \
	loms_sync_slave_db1 loms_master_db1 loms_migrate1 \
	loms_reset_stocks

run-dc-metrics:
	docker compose up --force-recreate --build prometheus grafana jaeger

run-dc-kafka:
	docker compose up --force-recreate --build kafka-ui kafka0

generate-all-mocks:
	cd cart && make generate-mocks && cd .. ;\
    cd loms && make generate-mocks && cd .. ;\
    cd common && make generate-mocks && cd ..

run-all-unit-tests: go-work-sync
	cd cart && make run-unit-tests ; cd .. && \
    cd loms && make run-unit-tests

ci-run-all-unit-tests:
	@ cd cart && make ci-run-unit-tests > /dev/null && \
 	printf "\nBegin cart report\n" && cat reports/report.xml && printf "\nEnd cart report\n" && \
 	printf "\nBegin cart coverage\n" && cat reports/coverage.xml && printf "\nEnd cart coverage\n" && \
 	cd .. && \
 	cd loms && make ci-run-unit-tests > /dev/null && \
 	printf "\nBegin loms report\n" && cat reports/report.xml && printf "\nEnd loms report\n"
#    printf "\nBegin loms coverage\n" && cat reports/coverage.xml && printf "\nEnd loms coverage\n"

run-all-unit-tests-coverage:
	cd cart && make run-unit-tests-coverage ; cd .. && \
	cd loms && make run-unit-tests-coverage

run-all-e2e-tests: build-all-for-linux
	@ \
 	cd cart && make run-e2e-tests ; cd .. && \
 	cd loms && make run-e2e-tests

run-all-sqlc-tests: build-all
	@ \
 	cd loms && make run-sqlc-tests

run-all-integration-tests: build-all
	@ \
 	cd loms && make run-sqlc-tests ; cd .. ;\
 	cd cart && make run-integration-tests


run-all-gocyclo-linters:
	cd cart && make run-gocyclo-linter ; cd .. && \
 	cd loms && make run-gocyclo-linter

run-all-gocognit-linters:
	cd cart && make run-gocognit-linter ; cd .. && \
 	cd loms && make run-gocognit-linter

run-all-benchmarks:
	cd cart && make run-benchmarks

go-work-sync:
	go work sync

generate-all-api:
	cd loms && make generate-api ;\
	cd .. ;\
	mkdir -p cart/api/loms/v1 && cp loms/api/loms/v1/loms.proto cart/api/loms/v1/loms.proto ;\
	cp -rf loms/vendor-proto/ cart/vendor-proto ;\
	cd cart && make generate-api && unlink api/loms/v1/loms.proto ;\
	cd .. ;\
	mkdir -p notifier/api/loms/v1 && cp loms/api/loms/v1/message.proto notifier/api/loms/v1/message.proto ;\
    cp -rf loms/vendor-proto/ notifier/vendor-proto ;\
    cd notifier && make generate-api && unlink api/loms/v1/message.proto ;\
    cd .. ;\


generate-all-sql:
	cd loms && make generate-sql && cd .. ;\
	cd transactionalbox && make generate-sql && cd ..

bin-deps-all:
	cd cart && make bin-deps && cd .. ;\
	cd loms && make bin-deps
