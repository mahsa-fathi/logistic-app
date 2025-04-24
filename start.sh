#!/bin/sh

cd /app

APP_TYPE=$1

if [ "$APP_TYPE" = "server" ]; then
    go build -o server ./cmd/http
    ./server
elif [ "$APP_TYPE" = "cron" ]; then
    go build -o cron ./cmd/cron
    ./cron
fi
