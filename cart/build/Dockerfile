FROM alpine:3.20 as run

COPY dist/app /usr/local/bin/app
RUN chmod a+x /usr/local/bin/app

CMD ["app"]
