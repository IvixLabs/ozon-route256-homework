FROM alpine:3.20 as base

RUN apk add curl

RUN curl -fsSL https://raw.githubusercontent.com/pressly/goose/master/install.sh | sh

COPY migrations /usr/local/src/app/migrations

COPY build/cmd.sh /usr/local/bin/cmd.sh
RUN chmod a+x /usr/local/bin/cmd.sh

CMD ["cmd.sh"]
