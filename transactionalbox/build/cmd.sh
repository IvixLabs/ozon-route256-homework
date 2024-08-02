#!/bin/sh

goose -dir /usr/local/src/app/migrations -table transactionalbox_goose_db_version up
