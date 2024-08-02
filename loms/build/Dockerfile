FROM alpine:3.20 as run

COPY dist/app /usr/local/bin/app
RUN chmod a+x /usr/local/bin/app

CMD ["app"]

FROM alpine:3.20 as base_migrate

RUN apk add curl

RUN curl -fsSL https://raw.githubusercontent.com/pressly/goose/master/install.sh | sh


FROM base_migrate as migrate

COPY migrations /usr/local/src/app/migrations

WORKDIR /usr/local/src/app/migrations

CMD ["goose", "up"]


FROM base_migrate as reset_stocks

COPY reset_stocks_migrations /usr/local/src/app/reset_stocks_migrations

WORKDIR /usr/local/src/app/reset_stocks_migrations

CMD ["goose", "-no-versioning", "up"]
