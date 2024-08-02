FROM golang:1.22.4-alpine3.20 as base

RUN apk add --no-cache make

WORKDIR /usr/src/app

RUN echo "version 1.0.1"

COPY go.work .
COPY go.work.sum .

COPY ./common/go.* ./common/

COPY ./launcher/go.* ./launcher/

COPY ./swagger/go.* ./swagger/

COPY ./logger/go.* ./logger/

COPY ./metrics/go.* ./metrics/

COPY ./pprof/go.* ./pprof/

COPY ./debugsrv/go.* ./debugsrv/

COPY ./transactionalbox/go.* ./transactionalbox/

COPY ./notifier/go.* ./notifier/
COPY ./notifier/Makefile ./notifier/

RUN mkdir -p ./cart
RUN mkdir -p ./cart/bin
RUN mkdir -p ./cart/reports
COPY ./cart/go.mod ./cart/
COPY ./cart/go.sum ./cart/
COPY ./cart/Makefile ./cart/

RUN mkdir -p ./loms
RUN mkdir -p ./loms/bin
RUN mkdir -p ./loms/reports
COPY ./loms/go.mod ./loms/
COPY ./loms/go.sum ./loms/
COPY ./loms/Makefile ./loms/

COPY Makefile .
COPY ./scripts/ ./scripts/

RUN make bin-deps-all
RUN make go-work-sync

COPY ./cart/build/ ./cart/build/
COPY ./cart/cmd/ ./cart/cmd/
COPY ./cart/internal/ ./cart/internal/
COPY ./cart/test/ ./cart/test/
COPY ./cart/.env.docker ./cart/

COPY ./loms/build/ ./loms/build/
COPY ./loms/cmd/ ./loms/cmd/
COPY ./loms/internal/ ./loms/internal/
COPY ./loms/test/ ./loms/test/
COPY ./loms/migrations/ ./loms/migrations/
COPY ./loms/test_migrations/ ./loms/test_migrations/
COPY ./loms/.env ./loms/

COPY ./notifier/build/ ./notifier/build/
COPY ./notifier/cmd/ ./notifier/cmd/
COPY ./notifier/internal/ ./notifier/internal/
COPY ./notifier/.env ./loms/

COPY ./logger/pkg/ ./logger/pkg/
COPY ./swagger/pkg/ ./swagger/pkg/
COPY ./pprof/pkg/ ./pprof/pkg/
COPY ./metrics/pkg/ ./metrics/pkg/
COPY ./common/pkg/ ./common/pkg/
COPY ./debugsrv/pkg/ ./debugsrv/pkg/
COPY ./transactionalbox/pkg/ ./transactionalbox/pkg/

