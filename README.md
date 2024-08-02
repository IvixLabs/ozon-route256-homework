# Homework of ozon route256 go middle course

# Project description

- Cart service call remote product service to fetch name of product and manage cart for user
- Loms service create order if stocks are available and send events to kafka
- Notifier service print events to console. There are 3 instances.
- There are 4 instances of postgres 2 pair of master-slave. Masters are used for sharding order info. Each master is used for writing and slave for reading.
- Kafka is used for communication between loms and notifiers
- Redis is used for cache remote product calls in cart


To run locally you need:
- make utility 
- Docker or podman 
- Clone project
- Run make run-all-dc from root folder


# Technologies and libraries
- Go lang (golang)
- goose
- pgx
- sqlc
- minimock
- protoc
- Postgresql
- Kafka
- Redis
- Jaeger
- Prometheus
- Grafana
- Pprof
- GRPC
- REST

# Structure

## Cart service
Contains REST API for manage user cart.

- list cart items ([example](./cart/examples/get_cart.http))
- add cart item ([example](./cart/examples/add_cart_item.http))
- remove cart item ([example](./cart/examples/remove_cart_item.http))
- remove whole cart ([example](./cart/examples/remove_cart.http))
- checkout cart  ([example](./cart/examples/checkout_cart.http))

## Loms service
Contains GRPC api for available stocks and pay/cancel order.
GRPC is available on localhost:50051 you can use postman for example to fetch GRPC schema

- OrderInfo - get order info(you can invoke directly after cart checkout)
- OrderCreat - create order if stocks are available(cart uses this method)
- OrderPay - pay for created order(you can invoke directly after cart checkout)
- OrderCancel - cancel created order(you can invoke directly after cart checkout)
- StockInfo  - get information about stock(cart uses this method)

Swagger at http://localhost:8084/swagger/

## Notifier service
Receives notifications about order from kafka and print to console.

## Grafana
http://localhost:3000

## Jaeger
http://localhost:16686

## Kafka ui
http://localhost:8090

## Unit tests
make run-all-unit-tests

## Integration tests
make run-all-integration-tests

## e2e tests
make run-all-e2e-tests

# Курсовой проект 12-го потока route256 go middle
- [Домашние задания](./docs/README.md)


